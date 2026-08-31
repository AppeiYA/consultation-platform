package domain

import (
	"strings"

	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

type Expertise struct {
	name string
}

var (
	ErrInvalidExpertise = custom_errors.BadException("invalid expertise name")
)

func NewExpertise(name string) (Expertise, error) {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) == 0 {
		return Expertise{}, ErrInvalidExpertise
	}
	return Expertise{name: trimmed}, nil
}

func (e Expertise) Name() string {
	return e.name
}

func (e Expertise) String() string {
	return e.name
}

func (e Expertise) Equals(other Expertise) bool {
	return strings.EqualFold(e.name, other.name)
}
