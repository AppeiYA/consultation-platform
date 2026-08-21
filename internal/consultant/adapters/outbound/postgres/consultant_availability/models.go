package consultantAvailabilityRepo

import (
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type ConsultantAvailability struct {
	ID string `db:"id"`
	ConsultantID string `db:"consultant_id"`
	DayOfWeek int `db:"day_of_week"`
	StartTime time.Time `db:"start_time"`
	EndTime time.Time `db:"end_time"`
	IsActive bool `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func ConsultantAvailabilityFromDomain(availability *domain.ConsultantAvailability) *ConsultantAvailability {
	startTime := availability.StartTime().Time(availability.DayOfWeek())
	endTime := availability.EndTime().Time(availability.DayOfWeek())

	return &ConsultantAvailability{
		ID: availability.ID(),
		ConsultantID: availability.ConsultantID(),
		DayOfWeek: int(availability.DayOfWeek()),
		StartTime: startTime,
		EndTime: endTime,
		IsActive: availability.IsActive(),
		CreatedAt: availability.CreatedAt(),
		UpdatedAt: availability.UpdatedAt(),
	}
}

func (a *ConsultantAvailability) ToDomain() (*domain.ConsultantAvailability, error) {
	startTime := domain.NewTimeOfDayFromTime(a.StartTime)
	endTime := domain.NewTimeOfDayFromTime(a.EndTime)

	return domain.ReconstitueConsultantAvailability(
		a.ID,
		a.ConsultantID,
		time.Weekday(a.DayOfWeek),
		startTime,
		endTime,
		a.IsActive,
		a.CreatedAt,
		a.UpdatedAt,
	)
}