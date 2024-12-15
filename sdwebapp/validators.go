package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdreflect"
	"github.com/gaorx/stardust6/sdvalidate"
	"github.com/labstack/echo/v4"
	"reflect"
)

const (
	akValidator         = "sdwebapp.validator"
	akValidationEnabled = "sdwebapp.validation_enabled"
)

var _ echo.Validator = ValidatorFunc(nil)

type ValidatorFunc func(v any) error

func (f ValidatorFunc) Validate(v any) error {
	return f(v)
}

func defaultValidate(v any) error {
	t := reflect.TypeOf(v)
	if sdreflect.IsStruct(t) || sdreflect.IsStructPtr(t) {
		return sdvalidate.Struct(v)
	}
	return nil
}

func validate(c echo.Context, v any) error {
	enabled := Get[bool](c, akValidationEnabled)
	if !enabled {
		return nil
	}
	validator := Get[echo.Validator](c, akValidator)
	if validator == nil {
		validator = c.Echo().Validator
	}
	if validator == nil {
		return nil
	}
	err := validator.Validate(v)
	if err != nil {
		if !sderr.Is(err, echo.ErrValidatorNotRegistered) {
			return err
		}
	}
	return nil
}
