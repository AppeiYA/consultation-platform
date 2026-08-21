package consultantAvailabilityRepo

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var saveAvailabilityQuery = `
INSERT INTO consultant_availabilities (id, consultant_id, day_of_week, start_time, end_time, created_at, updated_at)
VALUES (:id, :consultant_id, :day_of_week, :start_time, :end_time, :created_at, :updated_at)
`
func (a *AvailabilityRepository) SaveAvailability(ctx context.Context, availability *domain.ConsultantAvailability) error {
	executor := a.repository.Executor(ctx)
	model := ConsultantAvailabilityFromDomain(availability)

	_, err := executor.NamedExecContext(ctx, saveAvailabilityQuery, model)
	if err != nil {
		return err
	}
	return nil
}