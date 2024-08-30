package validation

import (
	"errors"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

var _ echo.Validator = (*CustomValidator)(nil)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// Custom error message for each field validation error.
func FormatValidationError(err error) map[string]interface{} {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return nil
	}

	errorsMap := make(map[string]interface{})
	for _, fieldErr := range validationErrors {
		fieldName := fieldErr.Field()
		key := strings.ToLower(fieldName) // Convert field name to lowercase
		switch fieldErr.Tag() {
		case "required":
			errorsMap[key] = fieldName + " is required"
		case "email":
			errorsMap[key] = "Invalid email format"
		case "min":
			errorsMap[key] = fieldName + " should be at least " + fieldErr.Param() + " characters long"
		// Add more case statements for different validation rules as needed
		default:
			errorsMap[key] = fieldName + " is not valid"
		}
	}
	return errorsMap
}
