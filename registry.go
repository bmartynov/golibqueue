package golibqueue

import "errors"

var (
	registry          = map[string]Factory{}
	ErrUnknownBackend = errors.New("golibqueue: unknown backend")
)

type Factory func(map[string]interface{}) (Queue, error)

func Create(backend string, rawCfg map[string]interface{}) (Queue, error) {
	f, ok := registry[backend]
	if !ok {
		return nil, ErrUnknownBackend
	}

	return f(rawCfg)
}

func Register(backend string, f Factory) {
	registry[backend] = f
}
