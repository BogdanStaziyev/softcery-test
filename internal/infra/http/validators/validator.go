package validators

import (
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"strings"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if err == nil {
		return nil
	}
	err = errors.New(strings.Replace(err.Error(), "\n", ", ", -1))
	return err
}
