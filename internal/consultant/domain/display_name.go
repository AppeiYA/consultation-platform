package domain

import (
	"regexp"
	"strings"
)

type DisplayName struct {
	value string
}

func NewDisplayName(displayName string) (DisplayName, error) {
	value := strings.TrimSpace(displayName)

	if len(value) < 6 || len(value) > 20 {
		return DisplayName{}, ErrInvalidDisplayNameLength
	}

	if !isAlphanumericRegex(value) {
		return DisplayName{}, ErrInvalidDisplayName
	}

	return DisplayName{value: value}, nil
}

func (d DisplayName) String() string {
	return d.value
}

func isAlphanumericRegex(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)
	return re.MatchString(s)
}