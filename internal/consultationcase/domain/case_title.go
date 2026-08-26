package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

type CaseTitle struct {
	value string
}

var (
	ErrCaseTitleEmpty = custom_errors.BadException("Case title is empty")
	ErrCaseTitleTooLong = custom_errors.BadException("Case title is too long")
)
func NewCaseTitle(value string) (*CaseTitle, error) {
	if len(value) == 0 {
		return nil, ErrCaseTitleEmpty
	}

	if len(value) > 100 {
		return nil, ErrCaseTitleTooLong
	}

	return &CaseTitle{value: value}, nil
}
func (c CaseTitle) String() string {return c.value}