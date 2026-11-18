package util

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates any struct and returns a friendly error message if validation fails.
func ValidateStruct(s any) error {
	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errMsgs []string
			for _, fieldErr := range validationErrors {
				errMsgs = append(errMsgs, fmt.Sprintf("%s is %s", strings.ToLower(fieldErr.Field()), fieldErr.Tag()))
			}
			return fmt.Errorf("validation failed: %s", strings.Join(errMsgs, ", "))
		}
		return err
	}
	return nil
}
