package nsq

import (
	"log"
	"sync"

	"github.com/nsqio/go-nsq"

	"github.com/bmartynov/golibqueue"
)

type Queue struct {
	sync.RWMutex
	logger    *log.Logger
	config    *Config
	producer  *nsq.Producer
	consumers map[string]*nsq.Consumer
}

func (q *Queue) buildConfig() (cfg *nsq.Config, err error) {
	cfg = nsq.NewConfig()
	err = cfg.Set("max_attempts", q.config.MaxAttempts)
	if err != nil {
		return cfg, err
	}

	err = cfg.Set("max_in_flight", q.config.MaxInFlight)
	if err != nil {
		return cfg, err
	}

	err = cfg.Set("default_requeue_delay", q.config.DefaultRequeueDelay.String())
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (q *Queue) createProducer() (p *nsq.Producer, err error) {
	cfg, err := q.buildConfig()
	if err != nil {
		return
	}

	p, err = nsq.NewProducer(q.config.NsqDS[0], cfg)
	if err != nil {
		return
	}

	p.SetLogger(q.logger, nsq.LogLevelError)

	return
}

func (q *Queue) Run() (err error) {
	q.producer, err = q.createProducer()
	if err != nil {
		return
	}

	return
}

func (q *Queue) Shutdown() error {
	q.producer.Stop()

	for name, _ := range q.consumers {
		q.consumers[name].Stop()
	}

	return nil
}

func (q *Queue) Subscribe(
	topic, channel string,
	handler golibqueue.Handler,
	opts ...golibqueue.Option,
) (err error) {

	so := &nsqOptions{}
	for _, opt := range opts {
		err = opt(so)
		if err != nil {
			return
		}
	}

	var wHandler = func(msg *nsq.Message) error {
		return handler.Handle(wrapMessage(msg))
	}

	var consumer *nsq.Consumer

	consumer, err = nsq.NewConsumer(topic, channel, so.cfg)
	if err != nil {
		return
	}

	consumer.SetLogger(
		q.logger,
		nsq.LogLevelError,
	)

	consumer.AddConcurrentHandlers(
		nsq.HandlerFunc(wHandler),
		q.config.ConsumerConcurrency,
	)

	err = consumer.ConnectToNSQDs(q.config.NsqDS)
	if err != nil {
		return
	}

	err = consumer.ConnectToNSQLookupds(q.config.NsqLookupDS)
	if err != nil {
		return
	}

	key := golibqueue.FormatTopicQueueName(topic, channel)

	q.Lock()
	q.consumers[key] = consumer
	q.Unlock()

	return nil
}

func (q *Queue) UnSubscribe(topic, channel string) (err error) {
	q.Lock()
	defer q.Unlock()

	key := golibqueue.FormatTopicQueueName(topic, channel)

	if consumer, ok := q.consumers[key]; ok {
		consumer.Stop()
		delete(q.consumers, key)
	}

	return
}

func (q *Queue) Publish(topic string, payload []byte) error {
	return q.producer.Publish(topic, payload)
}

func New(config *Config) golibqueue.Queue {
	q := Queue{
		config:    config,
		consumers: make(map[string]*nsq.Consumer),
	}

	return &q
}
