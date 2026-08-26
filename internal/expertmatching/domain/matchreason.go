package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

type MatchReason struct {
	factor MatchFactor
	detail string
}

var (
	ErrInvalidMatchReasonDetail = custom_errors.BadException("invalid match reason detail")
)

func NewMatchReason(factor MatchFactor, detail string) (MatchReason, error) {
	if len(detail) == 0 {
		return MatchReason{}, ErrInvalidMatchReasonDetail
	}
	return MatchReason{
		factor: factor,
		detail: detail,
	}, nil
}

func (mr MatchReason) Factor() MatchFactor {
	return mr.factor
}

func (mr MatchReason) Detail() string {
	return mr.detail
}