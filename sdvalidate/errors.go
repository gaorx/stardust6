package sdvalidate

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/go-playground/validator/v10"
)

func IsInvalidValidationErr(err error) bool {
	_, ok := ToInvalidValidationErr(err)
	return ok
}

func ToInvalidValidationErr(err error) (*validator.InvalidValidationError, bool) {
	if err == nil {
		return nil, false
	}
	if err1, ok := sderr.As[*validator.InvalidValidationError](sderr.Root(err)); ok {
		return err1, true
	} else {
		return nil, false
	}
}

func IsValidationErr(err error) bool {
	_, ok := ToValidationErrs(err)
	return ok
}

func ToValidationErrs(err error) ([]validator.FieldError, bool) {
	if err == nil {
		return nil, false
	}
	if fieldErrs, ok := sderr.As[validator.ValidationErrors](sderr.Root(err)); ok {
		return fieldErrs, true
	} else {
		return nil, false
	}
}
