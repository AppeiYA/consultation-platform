package consultantAvailabilityRepo

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var findAvailabilityByIDQuery = `
SELECT *
FROM consultant_availabilities
WHERE id = $1
`
func (a *AvailabilityRepository) FindAvailabilityByID(ctx context.Context, availabilityID string) (*domain.ConsultantAvailability, error) {
	executor := a.repository.Executor(ctx)

	var model ConsultantAvailability

	err := executor.GetContext(
		ctx,
		&model,
		findAvailabilityByIDQuery,
		availabilityID,
	)
	if err != nil {
		return nil, err
	}

	return model.ToDomain()
}