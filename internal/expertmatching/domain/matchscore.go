package domain

import (
	"math"

	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

var minVal, maxVal float64 = 0.0, 1.0

type MatchScore struct {
	value float64
}

var (
	ErrInvalidMatchScore = custom_errors.BadException("invalid match score")
)

func NewMatchScore(value float64) (MatchScore, error) {
	if math.IsNaN(value) || value < minVal || value > maxVal {
		return MatchScore{}, ErrInvalidMatchScore
	}
	return MatchScore{value: value}, nil
}

func (ms MatchScore) Value() float64 {return ms.value}