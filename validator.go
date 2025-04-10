package service

import (
	events "github.com/yeencloud/lib-events"
	"github.com/yeencloud/lib-shared/validation"
)

func NewValidator() (*validation.Validator, error) {
	validator, err := validation.NewValidator()

	if err != nil {
		return nil, err
	}

	err = validator.RegisterValidations(events.Validations())
	if err != nil {
		return nil, err
	}
	return validator, nil
}
