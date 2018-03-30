package stan

import (
	"log"
	"sync"

	"github.com/nats-io/go-nats-streaming"

	"github.com/bmartynov/golibqueue"
)

type stanQueue struct {
	sync.RWMutex
	logger    *log.Logger
	config    *Config
	conn      stan.Conn
	consumers map[string]stan.Subscription
}

func (q *stanQueue) Run() (err error) {
	q.conn, err = stan.Connect(
		q.config.ClusterID,
		q.config.ClientID,
		stan.NatsURL(q.config.servers()),
	)

	return
}

func (q *stanQueue) Shutdown() (err error) {
	return q.conn.Close()
}

func (q *stanQueue) Subscribe(
	topic, channel string,
	handler golibqueue.Handler,
	opts ...golibqueue.Option,
) (err error) {
	so := &stanOptions{}
	for _, opt := range opts {
		err = opt(so)
		if err != nil {
			return
		}
	}

	var subscription stan.Subscription

	var wHandler = func(msg *stan.Msg) {
		handler.Handle(wrapMessage(msg))
	}

	if channel == "" {
		subscription, err = q.conn.Subscribe(topic, wHandler, so.options...)
	} else {
		subscription, err = q.conn.QueueSubscribe(topic, channel, wHandler, so.options...)
	}
	if err != nil {
		return
	}

	key := golibqueue.FormatTopicQueueName(topic, channel)

	q.Lock()
	q.consumers[key] = subscription
	q.Unlock()

	return
}

func (q *stanQueue) UnSubscribe(topic, channel string) (err error) {
	q.Lock()
	defer q.Unlock()

	key := golibqueue.FormatTopicQueueName(topic, channel)

	if consumer, ok := q.consumers[key]; ok {
		err = consumer.Unsubscribe()
		delete(q.consumers, key)
	}

	return
}

func (q *stanQueue) Publish(topic string, payload []byte) (err error) {
	return q.conn.Publish(topic, payload)
}

func New(config *Config) *stanQueue {
	q := stanQueue{
		config:    config,
		consumers: make(map[string]stan.Subscription),
	}

	return &q
}
