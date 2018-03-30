package golibqueue

import (
	"context"
)

type HandlerFn func(Message) error

func (h HandlerFn) Handle(message Message) error {
	return h(message)
}

type Handler interface {
	Handle(Message) error
}

type Queue interface {
	Publish(string, []byte) error
	Subscribe(string, string, Handler, ...Option) error
}

type QueueWRequest interface {
	Request(ctx context.Context, subj string, data []byte, cb Handler) error
}

type RunAble interface {
	Run() error
}

type ShutdownAble interface {
	Shutdown() error
}
