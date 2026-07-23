package domain

import (
	"unicode"

	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

var (
	invalidPasswordLength           = custom_errors.BadException("password must be at least 8 characters long")
	invalidPasswordUppercase        = custom_errors.BadException("password must contain at least one uppercase letter")
	invalidPasswordLowercase        = custom_errors.BadException("password must contain at least one lowercase letter")
	invalidPasswordDigit            = custom_errors.BadException("password must contain at least one digit")
	invalidPasswordSpecialCharacter = custom_errors.BadException("password must contain at least one special character")
)

type Password struct {
	value string
}

func NewPassword(value string) (Password, error) {
	password := Password{
		value: value,
	}

	if err := password.validate(); err != nil {
		return Password{}, err
	}

	return password, nil
}

func (p Password) String() string {
	return p.value
}

func (p Password) validate() error {
	if len(p.value) < 8 {
		return invalidPasswordLength
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, r := range p.value {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r), unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	switch {
	case !hasUpper:
		return invalidPasswordUppercase
	case !hasLower:
		return invalidPasswordLowercase
	case !hasDigit:
		return invalidPasswordDigit
	case !hasSpecial:
		return invalidPasswordSpecialCharacter
	}

	return nil
}