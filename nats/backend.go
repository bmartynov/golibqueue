package nats

import (
	"sync"
	"log"
	"context"

	"github.com/nats-io/go-nats"

	"github.com/bmartynov/golibqueue"
)

type natsQueue struct {
	sync.RWMutex
	logger    *log.Logger
	config    *Config
	conn      *nats.Conn
	consumers map[string]*nats.Subscription
}

func (q *natsQueue) Run() (err error) {
	q.conn, err = q.config.opts().Connect()

	return
}

func (q *natsQueue) Shutdown() (err error) {
	q.conn.Close()

	return nil
}

func (q *natsQueue) Request(
	ctx context.Context,
	subj string,
	data []byte,
	cb golibqueue.Handler,
) (err error) {

	msg, err := q.conn.RequestWithContext(ctx, subj, data)
	if err != nil {
		return err
	}

	return cb.Handle(&message{msg})
}

func (q *natsQueue) Subscribe(
	topic, channel string,
	handler golibqueue.Handler,
	_ ...golibqueue.Option,
) (err error) {
	var subscription *nats.Subscription

	var wHandler = func(msg *nats.Msg) error {
		handler.Handle(wrapMessage(msg))
	}

	if channel == "" {
		subscription, err = q.conn.Subscribe(topic, wHandler)
	} else {
		subscription, err = q.conn.QueueSubscribe(topic, channel, wHandler)
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

func (q *natsQueue) UnSubscribe(topic, channel string) (err error) {
	q.Lock()
	defer q.Unlock()

	key := golibqueue.FormatTopicQueueName(topic, channel)

	if consumer, ok := q.consumers[key]; ok {
		err = consumer.Unsubscribe()
		delete(q.consumers, key)
	}

	return
}

func New(config *Config) *natsQueue {
	q := natsQueue{
		config:    config,
		consumers: make(map[string]*nats.Subscription),
	}

	return &q
}
