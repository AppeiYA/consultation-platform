package domain

import (
	"strings"

	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

const ExpertiseIDPrefix = "exp"

type Expertise struct {
	id           string
	consultantID string
	name         string
}

var (
	ErrInvalidExpertiseID   = custom_errors.BadException("invalid expertise ID")
	ErrInvalidConsultantID  = custom_errors.BadException("invalid consultant ID")
	ErrInvalidExpertiseName = custom_errors.BadException("expertise name cannot be empty")
	ErrDuplicateExpertise   = custom_errors.BadException("duplicate expertise for consultant")
)

func NewExpertise(id string, consultantID string, name string) (*Expertise, error) {
	if len(id) == 0 {
		return nil, ErrInvalidExpertiseID
	}
	if len(consultantID) == 0 {
		return nil, ErrInvalidConsultantID
	}
	trimmed := strings.TrimSpace(name)
	if len(trimmed) == 0 {
		return nil, ErrInvalidExpertiseName
	}

	return &Expertise{
		id:           id,
		consultantID: consultantID,
		name:         trimmed,
	}, nil
}

func (e *Expertise) ID() string {
	return e.id
}

func (e *Expertise) ConsultantID() string {
	return e.consultantID
}

func (e *Expertise) Name() string {
	return e.name
}
