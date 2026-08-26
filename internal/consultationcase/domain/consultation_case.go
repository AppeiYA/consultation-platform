package domain

import (
	"time"

	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

var CaseIDPrefix = "case"

var (
	ErrInvalidCaseID = custom_errors.BadException("Invalid case ID")
	ErrInvalidClientID = custom_errors.BadException("Invalid client ID")
	ErrInvalidStatusTransition = custom_errors.BadException("Invalid status transition")
)

type ConsultationCase struct {
	id string
	clientID string
	title CaseTitle
	description CaseDescription
	category CaseCategory
	status CaseStatus
	createdAt time.Time
	updatedAt time.Time
}

func NewConsultationCase(
	id string,
	clientID string,
	title CaseTitle,
	description CaseDescription,
	category CaseCategory,
	now time.Time,
) (*ConsultationCase, error) {

	if id == "" {
		return nil, ErrInvalidCaseID
	}

	if clientID == "" {
		return nil, ErrInvalidClientID
	}

	return &ConsultationCase{
		id:          id,
		clientID:    clientID,
		title:       title,
		description: description,
		category:    category,
		status:      CaseStatusDraft,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func ReconstituteConsultationCase(
	id string,
	clientID string,
	title string,
	description string,
	category string,
	status string,
	createdAt time.Time,
	updatedAt time.Time,
) (*ConsultationCase, error) {

	if id == "" {
		return nil, ErrInvalidCaseID
	}
	if clientID == "" {
		return nil, ErrInvalidClientID
	}

	caseTitle, err := NewCaseTitle(title)
	if err != nil {
		return nil, err
	}

	caseDescription, err := NewCaseDescription(description)
	if err != nil {
		return nil, err
	}

	caseCategory, err := NewCaseCategory(category)
	if err != nil {
		return nil, err
	}

	caseStatus, err := NewCaseStatus(status)
	if err != nil {
		return nil, err
	}

	return &ConsultationCase{
		id:          id,
		clientID:    clientID,
		title:       *caseTitle,
		description: *caseDescription,
		category:    *caseCategory,
		status:      caseStatus,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}, nil
}

func (c *ConsultationCase) ID() string {return c.id}
func (c *ConsultationCase) ClientID() string {return c.clientID}
func (c *ConsultationCase) Title() CaseTitle {return c.title}
func (c *ConsultationCase) Description() CaseDescription {return c.description}
func (c *ConsultationCase) Category() CaseCategory {return c.category}
func (c *ConsultationCase) Status() CaseStatus {return c.status}
func (c *ConsultationCase) CreatedAt() time.Time {return c.createdAt}
func (c *ConsultationCase) UpdatedAt() time.Time {return c.updatedAt}

func (c *ConsultationCase) UpdateTitle(newTitle CaseTitle) {
	c.title = newTitle
}
func (c *ConsultationCase) UpdateDescription(newDescription CaseDescription) {
	c.description = newDescription
}
func (c *ConsultationCase) UpdateCategory(newCategory CaseCategory) {
	c.category = newCategory
}

func (c *ConsultationCase) EndUpdating(now time.Time) {
	c.updatedAt = now
}

func (c *ConsultationCase) TransitionToStatus(newStatus CaseStatus, now time.Time) error {
	if !c.status.CanTransitionTo(newStatus) {
		return ErrInvalidStatusTransition
	}
	c.status = newStatus
	c.updatedAt = now
	return nil
}