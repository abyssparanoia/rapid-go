package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func Validate(s interface{}) error {
	if err := v.Struct(s); err != nil {
		return fmt.Errorf("failed to validate struct: %w", err)
	}
	return nil
}
