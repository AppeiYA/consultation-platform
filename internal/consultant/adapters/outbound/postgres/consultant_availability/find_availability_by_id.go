package consultantAvailabilityRepo

import (
	"context"
	"database/sql"
	"errors"

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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAvailabilityNotFound
		}
		return nil, err
	}

	return model.ToDomain()
}