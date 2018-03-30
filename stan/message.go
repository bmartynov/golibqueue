package stan

import (
	"strconv"

	"github.com/nats-io/go-nats-streaming"

	"github.com/bmartynov/golibqueue"
)

type message struct {
	msg *stan.Msg
}

func (m *message) Finish() error {
	return m.msg.Ack()
}

func (m *message) Id() string {
	return strconv.Itoa(int(m.msg.Sequence))
}

func (m *message) Timestamp() int64 {
	return m.msg.Timestamp
}

func (m *message) Payload() []byte {
	return m.msg.Data
}

func (m *message) IsRedelivered() bool {
	return m.msg.Redelivered
}

func (m *message) Reply() string {
	return m.msg.Reply
}

func (m *message) Subject() string {
	return m.msg.Subject
}

func wrapMessage(msg *stan.Msg) golibqueue.Message {
	return &message{msg}
}
