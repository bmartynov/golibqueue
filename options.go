package golibqueue

type Options interface {
	Set(string, ...interface{}) error
}

type Option func(Options) error
