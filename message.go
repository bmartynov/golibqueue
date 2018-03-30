package golibqueue

import "time"

type Message interface {
	Payload() []byte
}

type WFinish interface {
	Finish() error
}

type WId interface {
	Id() string
}

type WTimestamp interface {
	Timestamp() int64
}

type WRequeue interface {
	Requeue(time.Duration)
}

type WIsRedelivered interface {
	IsRedelivered() bool
}

type WIsResponded interface {
	IsResponded() bool
}

type WReply interface {
	Reply() string
}

type WSubject interface {
	Subject() string
}
