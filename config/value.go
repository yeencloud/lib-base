package config

import (
	"fmt"
	"strconv"

	"github.com/fatih/structs"

	"github.com/yeencloud/lib-shared"
)

type valueHandler func(field *structs.Field, value string) error

func handleString(field *structs.Field, value string) error {
	return field.Set(value)
}

func handleSecret(field *structs.Field, value string) error {
	return field.Set(shared.NewSecret(value))
}

func wrapConversionError(err1 error, err2 error) error {
	return fmt.Errorf("%w: %w", err1, err2)
}

func handleInt(field *structs.Field, value string) error {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return wrapConversionError(UnsupportedValueForConversionError{
			Value:    value,
			FromType: "string",
			ToType:   "int",
		}, err)
	}

	err = field.Set(intValue)
	if err != nil {
		return err
	}

	return nil
}

func handleBool(field *structs.Field, value string) error {
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return wrapConversionError(UnsupportedValueForConversionError{
			Value:    value,
			FromType: "string",
			ToType:   "bool",
		}, err)
	}

	err = field.Set(boolValue)
	if err != nil {
		return err
	}

	return nil
}
