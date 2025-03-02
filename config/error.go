package config

import (
	"fmt"
)

type MissingConfigValueError struct {
	Key string
}

func (e MissingConfigValueError) Error() string {
	return "missing config value for key " + e.Key
}

type UnsupportedConfigTypeError struct {
	Type string

	Variable       string
	AvailableTypes []string
}

func (e UnsupportedConfigTypeError) Error() string {
	return "unsupported type " + e.Type + " for variable " + e.Variable
}

type UnsupportedValueForConversionError struct {
	Value string

	FromType string
	ToType   string
}

func (e UnsupportedValueForConversionError) Error() string {
	return fmt.Sprintf("failed to convert %s for conversion from %s to %s", e.Value, e.FromType, e.ToType)
}
