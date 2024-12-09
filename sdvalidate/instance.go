package sdvalidate

import (
	"github.com/go-playground/validator/v10"
)

var defaultValidate = validator.New(validator.WithRequiredStructEnabled())

func Default() *validator.Validate {
	return defaultValidate
}

func SetDefault(v *validator.Validate) {
	defaultValidate = v
}
