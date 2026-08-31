package domain

import (
	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

const (
	DefaultMaxCandidatePoolSize = 500
)

type CandidatePool struct {
	candidates []CandidateProfile
}

var (
	ErrCandidatePoolExceeded = custom_errors.BadException("candidate pool size exceeds maximum limit")
	ErrDuplicateCandidateInPool = custom_errors.BadException("duplicate consultant in candidate pool")
)

func NewCandidatePool(candidates []CandidateProfile, maxCapacity int) (CandidatePool, error) {
	if maxCapacity <= 0 {
		maxCapacity = DefaultMaxCandidatePoolSize
	}

	if len(candidates) > maxCapacity {
		return CandidatePool{}, ErrCandidatePoolExceeded
	}

	seen := make(map[string]struct{}, len(candidates))
	cloned := make([]CandidateProfile, 0, len(candidates))

	for _, c := range candidates {
		if _, exists := seen[c.ConsultantID()]; exists {
			return CandidatePool{}, ErrDuplicateCandidateInPool
		}
		seen[c.ConsultantID()] = struct{}{}
		cloned = append(cloned, c)
	}

	return CandidatePool{
		candidates: cloned,
	}, nil
}

func (cp CandidatePool) Candidates() []CandidateProfile {
	if len(cp.candidates) == 0 {
		return []CandidateProfile{}
	}
	cloned := make([]CandidateProfile, len(cp.candidates))
	copy(cloned, cp.candidates)
	return cloned
}

func (cp CandidatePool) Size() int {
	return len(cp.candidates)
}

func (cp CandidatePool) IsEmpty() bool {
	return len(cp.candidates) == 0
}
