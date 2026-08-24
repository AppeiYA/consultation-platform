package domain

import (
	"time"

	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

var ConsultantAvailabilityIDPrefix = "conav"

var (
	ErrInvalidTimeRange = custom_errors.BadException("invalid time range")
	ErrAvailabilityOverlap = custom_errors.ConflictError("availability overlaps with existing availability")
	ErrAvailabilityNotFound = custom_errors.NotFoundError("availability not found")
	ErrAvailabilityAlreadyDeactivated = custom_errors.ConflictError("availability is already deactivated")
	ErrAvailabilityAlreadyActivated = custom_errors.ConflictError("availability is already activated")
)

type ConsultantAvailability struct {
	id string
	consultantID string
	dayOfWeek time.Weekday
	startTime TimeOfDay
	endTime TimeOfDay
	isActive bool
	createdAt time.Time
	updatedAt time.Time
}

func NewConsultantAvailability(
	id string, 
	consultantID string, 
	dayOfWeek time.Weekday, 
	startTime TimeOfDay, 
	endTime TimeOfDay,
	now time.Time,
) (*ConsultantAvailability, error) {
	if !startTime.Before(endTime) {
		return nil, ErrInvalidTimeRange
	}
	return &ConsultantAvailability{
		id: id,
		consultantID: consultantID,
		dayOfWeek: dayOfWeek,
		startTime: startTime,
		endTime: endTime,
		isActive: true,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func ReconstitueConsultantAvailability(
	id string,
	consultantID string,
	dayOfWeek time.Weekday,
	startTime TimeOfDay,
	endTime TimeOfDay,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) (*ConsultantAvailability, error) {
	if !startTime.Before(endTime) {
		return nil, ErrInvalidTimeRange
	}
	return &ConsultantAvailability{
		id: id,
		consultantID: consultantID,
		dayOfWeek: dayOfWeek,
		startTime: startTime,
		endTime: endTime,
		isActive: isActive,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (ca *ConsultantAvailability) ID() string {
	return ca.id
}

func (ca *ConsultantAvailability) ConsultantID() string {
	return ca.consultantID
}

func (ca *ConsultantAvailability) DayOfWeek() time.Weekday {
	return ca.dayOfWeek
}

func (ca *ConsultantAvailability) StartTime() TimeOfDay {
	return ca.startTime
}

func (ca *ConsultantAvailability) EndTime() TimeOfDay {
	return ca.endTime
}

func (ca *ConsultantAvailability) IsActive() bool {
	return ca.isActive
}

func (ca *ConsultantAvailability) CreatedAt() time.Time {
	return ca.createdAt
}

func (ca *ConsultantAvailability) UpdatedAt() time.Time {
	return ca.updatedAt
}

func (a *ConsultantAvailability) Activate(now time.Time) {
	a.isActive = true
	a.updatedAt = now
}

func (a *ConsultantAvailability) Deactivate(now time.Time) {
	a.isActive = false
	a.updatedAt = now
}

func (a *ConsultantAvailability) Overlaps(
	start TimeOfDay,
	end TimeOfDay,
) bool {
	if !start.Before(end) {
		return false
	}

	return start.Before(a.endTime) && end.After(a.startTime)
}

func (a *ConsultantAvailability) Update(
	dayOfWeek time.Weekday,
	startTime TimeOfDay,
	endTime TimeOfDay,
	now time.Time,
) error {
	if !startTime.Before(endTime) {
		return ErrInvalidTimeRange
	}

	a.dayOfWeek = dayOfWeek
	a.startTime = startTime
	a.endTime = endTime
	a.updatedAt = now

	return nil
}

