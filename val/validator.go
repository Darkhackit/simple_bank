package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile("^[a-z0-9_]+$").MatchString
	isValidName     = regexp.MustCompile("^[a-zA-Z\\s]+$").MatchString
)

func ValidateString(value string, minLength, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("value length must be between %d and %d", minLength, maxLength)
	}
	return nil
}
func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("invalid username")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("invalid email")
	}
	return nil
}
func ValidateName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidName(value) {
		return fmt.Errorf("invalid name")
	}
	return nil
}
