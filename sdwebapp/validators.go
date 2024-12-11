package sdwebapp

import (
	"github.com/gaorx/stardust6/sdreflect"
	"github.com/gaorx/stardust6/sdvalidate"
	"github.com/labstack/echo/v4"
	"reflect"
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
