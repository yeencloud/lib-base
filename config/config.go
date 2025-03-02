package config

import (
	"reflect"

	"github.com/fatih/structs"
	"github.com/rs/zerolog/log"

	"github.com/yeencloud/lib-shared"

	"github.com/yeencloud/lib-base/config/source"
	"github.com/yeencloud/lib-base/depinjection"
)

type Config struct {
	dig depinjection.DependencyInjection

	sourceInterface source.ConfigInterface

	typeMap map[reflect.Kind]valueHandler
}

func NewConfig(dig depinjection.DependencyInjection, sourceInterface source.ConfigInterface) *Config {
	var secret shared.Secret

	config := &Config{
		dig:             dig,
		sourceInterface: sourceInterface,
		typeMap: map[reflect.Kind]valueHandler{
			reflect.TypeOf(secret).Kind(): handleSecret,
			reflect.String:                handleString,
			reflect.Int:                   handleInt,
			reflect.Bool:                  handleBool,
		},
	}

	return config
}

func (cfg *Config) AvailableTypes() []string {
	types := make([]string, 0, len(cfg.typeMap))
	for k := range cfg.typeMap {
		types = append(types, k.String())
	}
	return types
}

func RegisterConfig[obj any](config *Config) error {
	err := config.dig.Provide(func() (*obj, error) {
		var object obj

		s := structs.New(&object)

		for _, field := range s.Fields() {
			configKey := field.Tag("config")
			if configKey == "" {
				continue
			}

			value, err := config.sourceInterface.ReadString(configKey)
			if err != nil {
				return nil, err
			}

			defaultValue := field.Tag("default")
			if value == "" {
				value = defaultValue
			}

			handler, found := config.typeMap[field.Kind()]
			if !found {
				return nil, UnsupportedConfigTypeError{Type: field.Kind().String(), Variable: field.Name(), AvailableTypes: config.AvailableTypes()}
			}

			err = handler(field, value)
			log.Debug().Msgf("%s = %v", configKey, field.Value())
			if err != nil {
				return nil, err
			}
		}
		return &object, nil
	})
	return err
}
