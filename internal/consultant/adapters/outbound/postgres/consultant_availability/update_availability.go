package consultantAvailabilityRepo

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)
var updateAvailabilityQuery = `
UPDATE consultant_availabilities
SET
	day_of_week = :day_of_week,
	start_time = :start_time,
	end_time = :end_time,
	updated_at = :updated_at
WHERE id = :id
`
func (a *AvailabilityRepository) UpdateAvailability(ctx context.Context, availability *domain.ConsultantAvailability) error {
	executor := a.repository.Executor(ctx)
	model := ConsultantAvailabilityFromDomain(availability)

	_, err := executor.NamedExecContext(ctx, updateAvailabilityQuery, model)
	if err != nil {
		return err
	}
	return nil
}