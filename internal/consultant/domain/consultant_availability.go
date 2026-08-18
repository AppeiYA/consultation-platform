package domain

import (
	"time"

	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

var ConsultantAvailabilityIDPrefix = "conav"

var (
	ErrInvalidTimeRange = custom_errors.BadException("invalid time range")
	ErrAvailabilityOverlap = custom_errors.ConflictError("availability overlaps with existing availability")
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
	clock shared_outbound.Clock,
) (*ConsultantAvailability, error) {
	if !startTime.Before(endTime) {
		return nil, ErrInvalidTimeRange
	}
	now := clock.Now()
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

func (a *ConsultantAvailability) Activate(clock shared_outbound.Clock) {
	a.isActive = true
	a.updatedAt = clock.Now()
}

func (a *ConsultantAvailability) Deactivate(clock shared_outbound.Clock) {
	a.isActive = false
	a.updatedAt = clock.Now()
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

