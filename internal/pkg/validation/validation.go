package validation

import (
	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func Validate(s interface{}) error {
	if err := v.Struct(s); err != nil {
		return err
	}
	return nil
}
