package nats

import (
	"time"

	"github.com/nats-io/go-nats"
)

type Config struct {
	Addresses        []string       `yaml:"addresses" json:"addresses" mapstructure:"addresses"`
	ClientId         string         `yaml:"client_id" json:"client_id" mapstructure:"client_id"`
	AllowReconnect   *bool          `yaml:"allow_reconnects" json:"allow_reconnects" mapstructure:"allow_reconnects"`
	MaxReconnect     *int           `yaml:"max_reconnects" json:"max_reconnects" mapstructure:"max_reconnects"`
	ReconnectWait    *time.Duration `yaml:"reconnects_waits" json:"reconnects_waits" mapstructure:"reconnects_waits"`
	Timeout          *time.Duration `yaml:"timeout" json:"timeout" mapstructure:"timeout"`
	FlusherTimeout   *time.Duration `yaml:"flusher_timeout" json:"flusher_timeout" mapstructure:"flusher_timeout"`
	PingInterval     *time.Duration `yaml:"ping_interval" json:"ping_interval" mapstructure:"ping_interval"`
	MaxPingsOut      *int           `yaml:"max_ping_out" json:"max_ping_out" mapstructure:"max_ping_out"`
	ReconnectBufSize *int           `yaml:"reconnect_buf_size" json:"reconnect_buf_size" mapstructure:"reconnect_buf_size"`
	SubChanLen       *int           `yaml:"sub_chan_len" json:"sub_chan_len" mapstructure:"sub_chan_len"`
	User             *string        `yaml:"user" json:"user" mapstructure:"user"`
	Password         *string        `yaml:"password" json:"password" mapstructure:"password"`
	Token            *string        `yaml:"token" json:"token" mapstructure:"token"`
}

func (q *Config) opts() *nats.Options {
	opts := nats.GetDefaultOptions()

	opts.Servers = q.Addresses

	if q.AllowReconnect != nil && !(*q.AllowReconnect) {
		opts.AllowReconnect = false
	}
	if q.MaxReconnect != nil {
		opts.MaxReconnect = *q.MaxReconnect
	}
	if q.ReconnectWait != nil {
		opts.ReconnectWait = *q.ReconnectWait
	}
	if q.Timeout != nil {
		opts.Timeout = *q.Timeout
	}
	if q.FlusherTimeout != nil {
		opts.FlusherTimeout = *q.FlusherTimeout
	}
	if q.PingInterval != nil {
		opts.PingInterval = *q.PingInterval
	}
	if q.MaxPingsOut != nil {
		opts.MaxPingsOut = *q.MaxPingsOut
	}
	if q.ReconnectBufSize != nil {
		opts.ReconnectBufSize = *q.ReconnectBufSize
	}
	if q.SubChanLen != nil {
		opts.SubChanLen = *q.SubChanLen
	}
	if q.User != nil && q.Password != nil {
		opts.User = *q.User
		opts.Password = *q.Password
	}
	if q.Token != nil {
		opts.Token = *q.Token
	}

	return &opts
}

func (c *Config) Validate() error {
	return nil
}

func DefaultConfig() *Config {
	return &Config{}
}

type natsOptions struct {
	options []nats.Option
}

func (o *natsOptions) addOption(opt nats.Option) {
	o.options = append(o.options, opt)
}

func (o *natsQueue) Set(option string, value ...interface{}) (err error) {

}