package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

type CaseDescription struct {
	value string
}

var (
	ErrCaseDescriptionEmpty = custom_errors.BadException("Case description is empty")
	ErrCaseDescriptionTooLong = custom_errors.BadException("Case description is too long")
)
func NewCaseDescription(value string) (*CaseDescription, error) {
	if len(value) == 0 {
		return nil, ErrCaseDescriptionEmpty
	}

	if len(value) > 500 {
		return nil, ErrCaseDescriptionTooLong
	}

	return &CaseDescription{value: value}, nil
}

func (c CaseDescription) String() string {return c.value}