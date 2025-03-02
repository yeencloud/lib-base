package depinjection

type DependencyInjection interface {
	Invoke(function interface{}) error
	Provide(constructor interface{}) error
}

func NewDI() DependencyInjection {
	return NewDigWrapper()
}
