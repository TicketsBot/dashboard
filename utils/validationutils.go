package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func FormatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "max":
		if err.Type() == reflect.TypeOf("") {
			return fmt.Sprintf("Field \"%s\" cannot exceed %s characters in length", err.Field(), err.Param())
		} else {
			return fmt.Sprintf("Field \"%s\" cannot be greater than %s", err.Field(), err.Param())
		}
	case "min":
		if err.Type() == reflect.TypeOf("") {
			return fmt.Sprintf("Field \"%s\" must be at least %s characters in length", err.Field(), err.Param())
		} else {
			return fmt.Sprintf("Field \"%s\" cannot be less than %s", err.Field(), err.Param())
		}
	case "required":
		return fmt.Sprintf("Field \"%s\" is required", err.Field())
	default:
		return err.Error()
	}
}

func FormatValidationErrors(errors validator.ValidationErrors) string {
	var formatted string
	for _, err := range errors {
		formatted += FormatValidationError(err) + "\n"
	}

	formatted = strings.TrimSuffix(formatted, "\n")
	return formatted
}
