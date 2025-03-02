package depinjection

import (
	"go.uber.org/dig"
)

type DigWrapper struct {
	dig *dig.Container
}

func (d DigWrapper) wrapError(err error) error {
	if err != nil {
		return dig.RootCause(err)
	}
	return nil
}

func (d DigWrapper) Invoke(function interface{}) error {
	return d.wrapError(d.dig.Invoke(function))
}

func (d DigWrapper) Provide(constructor interface{}) error {
	return d.wrapError(d.dig.Provide(constructor))
}

func NewDigWrapper() DigWrapper {
	return DigWrapper{dig: dig.New()}
}
