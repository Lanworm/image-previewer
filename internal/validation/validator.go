package validation

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Validate(s interface{}) error {
	err := validate.Struct(s)

	return err
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}
