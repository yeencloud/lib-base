package config

type Environment struct {
	Name string `config:"ENV" default:"development"`
}

func (e Environment) IsProduction() bool {
	return e.Name == "production" || e.Name == "prod"
}
