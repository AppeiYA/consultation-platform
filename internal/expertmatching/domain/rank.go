package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

type Rank struct {
	value int
}

var (
	ErrInvalidRank = custom_errors.BadException("rank must be greater than or equal to 1")
)

func NewRank(value int) (Rank, error) {
	if value < 1 {
		return Rank{}, ErrInvalidRank
	}
	return Rank{value: value}, nil
}

func (r Rank) Value() int {
	return r.value
}
