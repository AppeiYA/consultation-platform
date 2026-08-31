package domain

import (
	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

type CandidateProfile struct {
	consultantID    string
	category        MatchingCategory
	profession      string
	expertise       []Expertise
	yearsExperience int
	bio             string
}

var (
	ErrInvalidYearsExperience = custom_errors.BadException("years of experience cannot be negative")
)

func NewCandidateProfile(
	consultantID string,
	category MatchingCategory,
	profession string,
	expertise []Expertise,
	yearsExperience int,
	bio string,
) (CandidateProfile, error) {
	if len(consultantID) == 0 {
		return CandidateProfile{}, ErrInvalidConsultantID
	}
	if yearsExperience < 0 {
		return CandidateProfile{}, ErrInvalidYearsExperience
	}

	var clonedExpertise []Expertise
	if len(expertise) > 0 {
		clonedExpertise = make([]Expertise, len(expertise))
		copy(clonedExpertise, expertise)
	}

	return CandidateProfile{
		consultantID:    consultantID,
		category:        category,
		profession:      profession,
		expertise:       clonedExpertise,
		yearsExperience: yearsExperience,
		bio:             bio,
	}, nil
}

func (cp CandidateProfile) ConsultantID() string {
	return cp.consultantID
}

func (cp CandidateProfile) Category() MatchingCategory {
	return cp.category
}

func (cp CandidateProfile) Profession() string {
	return cp.profession
}

func (cp CandidateProfile) Expertise() []Expertise {
	if len(cp.expertise) == 0 {
		return nil
	}
	cloned := make([]Expertise, len(cp.expertise))
	copy(cloned, cp.expertise)
	return cloned
}

func (cp CandidateProfile) YearsExperience() int {
	return cp.yearsExperience
}

func (cp CandidateProfile) Bio() string {
	return cp.bio
}
