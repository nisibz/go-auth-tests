package http

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) []string {
	var errMessages []string
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		errMessages = make([]string, 0, len(validationErrs))
		for _, fieldErr := range validationErrs {
			var msg string
			switch fieldErr.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", fieldErr.Field())
			case "email":
				msg = fmt.Sprintf("%s must be a valid email address", fieldErr.Field())
			case "min":
				msg = fmt.Sprintf("%s must be at least %s characters long", fieldErr.Field(), fieldErr.Param())
			default:
				msg = fmt.Sprintf("%s is not valid", fieldErr.Field())
			}
			errMessages = append(errMessages, msg)
		}
	}
	return errMessages
}
