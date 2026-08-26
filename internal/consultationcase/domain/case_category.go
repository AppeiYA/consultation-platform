package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

type CaseCategory struct {
	value string
}

var (
	ErrCaseCategoryEmpty = custom_errors.BadException("Case category is empty")
)
func NewCaseCategory(value string) (*CaseCategory, error) {
	if len(value) == 0 {
		return nil, ErrCaseCategoryEmpty
	}
	return &CaseCategory{value: value}, nil
}

func (c CaseCategory) String() string {return c.value}