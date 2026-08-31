package domain

import (
	"strings"

	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

type MatchingCategory struct {
	value string
}

var (
	ErrInvalidMatchingCategory = custom_errors.BadException("invalid matching category")
)

func NewMatchingCategory(value string) (MatchingCategory, error) {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) == 0 {
		return MatchingCategory{}, ErrInvalidMatchingCategory
	}
	return MatchingCategory{value: trimmed}, nil
}

func (c MatchingCategory) Value() string {
	return c.value
}

func (c MatchingCategory) String() string {
	return c.value
}

func (c MatchingCategory) Equals(other MatchingCategory) bool {
	return strings.EqualFold(c.value, other.value)
}
