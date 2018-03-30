package stan

import (
	"time"
	"strings"
	"fmt"

	"github.com/nats-io/go-nats-streaming"

	"github.com/bmartynov/libqueue"
)

type Config struct {
	ClientID              string     `mapstructure:"client_id"`
	ClusterID             string     `mapstructure:"cluster_id"`
	Addresses             []string   `mapstructure:"addresses"`
	DurableName           *string    `mapstructure:"durable_name"`
	AckWait               *time.Time `mapstructure:"ack_wait"`
	StartWithLastReceived *bool      `mapstructure:"start_with_last_received"`
	ReceiveAllAvailable   *bool      `mapstructure:"receive_all_available"`
	StartSequence         *uint64    `mapstructure:"start_sequence"`
	StartTime             *time.Time `mapstructure:"start_time"`
	ManualAckMode         *bool      `mapstructure:"manual_ack_mode"`
}

func (c *Config) servers() string {
	return strings.Join(c.Addresses, ",")
}

func (c *Config) Validate() error {
	return nil
}

func DefaultConfig() *Config {
	return &Config{}
}

//options
type stanOptions struct {
	options []stan.SubscriptionOption
}

func (o *stanOptions) addOption(opt stan.SubscriptionOption) {
	o.options = append(o.options, opt)
}

func (o *stanOptions) Set(option string, value ...interface{}) (err error) {
	switch option {
	case "durable_name":
		o.addOption(stan.DurableName(value[0].(string)))
	case "ack_wait":
		o.addOption(stan.AckWait(value[0].(time.Duration)))
	case "start_with_last_received":
		o.addOption(stan.StartWithLastReceived())
	case "receive_all_available":
		o.addOption(stan.DeliverAllAvailable())
	case "manual_ack_mode":
		o.addOption(stan.SetManualAckMode())
	case "start_sequence":
		o.addOption(stan.StartAtSequence(uint64(value[0].(int))))
	case "start_time":
		o.addOption(stan.StartAtTime(value[0].(time.Time)))
	case "manual_acks":
		o.addOption(stan.SetManualAckMode())
	default:
		err = fmt.Errorf("unknown option `%s` for stan", option)
	}
	return
}

func NewOptions() queue.Options {
	return &stanOptions{make([]stan.SubscriptionOption, 0)}
}
