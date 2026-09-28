package utils

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validateInstance *validator.Validate
	once             sync.Once
)

func getValidator() *validator.Validate {
	once.Do(func() {
		validateInstance = validator.New()
	})
	return validateInstance
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateErrors validates struct fields and returns structured error messages
func ValidateErrors(data any) ([]FieldError, error) {
	v := getValidator()
	err := v.Struct(data)
	if err == nil {
		return nil, nil
	}

	var errors []FieldError
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			var message string
			switch e.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", e.Field())
			case "email":
				message = "Please enter a valid email format"
			case "gte":
				message = fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
			case "gt":
				message = fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param())
			case "lte":
				message = fmt.Sprintf("%s must be less than or equal to %s", e.Field(), e.Param())
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param())
			case "max":
				message = fmt.Sprintf("%s cannot exceed %s characters", e.Field(), e.Param())
			default:
				message = fmt.Sprintf("%s failed validation on '%s'", e.Field(), e.Tag())
			}

			errors = append(errors, FieldError{
				Field:   e.Field(),
				Message: message,
			})
		}
		return errors, err
	}

	return nil, err
}

// ValidateInput helper returning string error for backwards compatibility
func ValidateInput(data interface{}) (string, error) {
	fieldErrors, err := ValidateErrors(data)
	if err != nil && len(fieldErrors) > 0 {
		return fieldErrors[0].Message, err
	}
	return "", err
}
