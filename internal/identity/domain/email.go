package domain

import (
	"net/mail"

	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

type Email struct {
	value string
}

var (
	invalidEmail = custom_errors.BadException("invalid email address")
)

func (e Email) Validate() (bool, error) {
	address, err := mail.ParseAddress(e.value)
	if err != nil {
		return false, invalidEmail
	}
	if address.Address != e.value {
		return false, nil
	}
	return true, nil
}

func (e Email) String() string {
	return e.value
}

func NewEmail(value string) (Email, error) {
	email := Email{value: value}
	if ok, err := email.Validate(); !ok {
		return Email{}, err
	}
	return email, nil
}
