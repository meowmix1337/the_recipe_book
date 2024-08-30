package validation

import (
	"unicode"
)

const minPasswordLength = 8

// ValidatePassword checks for password strength.
func ValidatePassword(password string) []string {
	var errors []string

	if len(password) < minPasswordLength {
		errors = append(errors, "password must be at least 8 characters long")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errors = append(errors, "password must contain at least one uppercase letter")
	}
	if !hasLower {
		errors = append(errors, "password must contain at least one lowercase letter")
	}
	if !hasNumber {
		errors = append(errors, "password must contain at least one number")
	}
	if !hasSpecial {
		errors = append(errors, "password must contain at least one special character")
	}

	return errors
}
