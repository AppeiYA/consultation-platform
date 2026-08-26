package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

type RankingVersion struct {
    value string
}

var (
	ErrInvalidRankingVersion = custom_errors.BadException("invalid ranking version")
)

func NewRankingVersion(value string) (RankingVersion, error) {
	if len(value) == 0 {
		return RankingVersion{}, ErrInvalidRankingVersion
	}
	return RankingVersion{value: value}, nil
}

func (rv RankingVersion) Value() string {
	return rv.value
}