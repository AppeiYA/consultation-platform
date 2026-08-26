package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

type MatchFactor string

const (
	MatchFactorExpertise MatchFactor = "expertise"
	MatchFactorAvailability MatchFactor = "availability"
	MatchFactorExperience MatchFactor = "experience"
	MatchFactorCategory MatchFactor = "category"
	MatchFactorLanguage MatchFactor = "language"
	MatchFactorVerification MatchFactor = "verification"
)

var (
	ErrInvalidMatchFactor = custom_errors.BadException("invalid match factor")
)

func (mf MatchFactor) IsValid() bool {
	switch mf {
	case MatchFactorExpertise, MatchFactorAvailability, MatchFactorExperience, MatchFactorCategory, MatchFactorLanguage, MatchFactorVerification:
		return true
	default:
		return false
	}
}

func NewMatchFactorFromString(s string) (MatchFactor, error) {
	mf := MatchFactor(s)
	if !mf.IsValid() {
		return "", ErrInvalidMatchFactor
	}
	return mf, nil
}

func (mf MatchFactor) String() string {
	return string(mf)
}