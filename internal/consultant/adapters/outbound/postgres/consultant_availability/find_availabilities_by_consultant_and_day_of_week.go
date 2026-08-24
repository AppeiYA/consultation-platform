package consultantAvailabilityRepo

import (
	"context"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var findAvailabilitiesByConsultantIDAndDayOfWeekQuery = `
SELECT id, consultant_id, day_of_week, start_time, end_time, is_active, created_at, updated_at
FROM consultant_availabilities
WHERE consultant_id = $1 AND day_of_week = $2
`
func (a *AvailabilityRepository) FindAvailabilitiesByConsultantIDAndDayOfWeek(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
	executor := a.repository.Executor(ctx)

	var models []ConsultantAvailability

	err := executor.SelectContext(
		ctx,
		&models,
		findAvailabilitiesByConsultantIDAndDayOfWeekQuery,
		consultantID,
		int(dayOfWeek),
	)
	if err != nil {
		return nil, err
	}

	availabilities := make([]*domain.ConsultantAvailability, 0, len(models))

	for _, model := range models {
		availability, err := model.ToDomain()
		if err != nil {
			return nil, err
		}

		availabilities = append(availabilities, availability)
	}

	return availabilities, nil
}