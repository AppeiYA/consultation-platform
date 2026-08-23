package consultantAvailabilityRepo

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var findAvailabilitiesByConsultantIDQuery = `
SELECT * FROM consultant_availabilities WHERE consultant_id = $1 AND is_active = true
`
func (a *AvailabilityRepository) FindAvailabilitiesByConsultantID(
    ctx context.Context,
    consultantID string,
) ([]*domain.ConsultantAvailability, error) {
    executor := a.repository.Executor(ctx)

    var models []ConsultantAvailability

    err := executor.SelectContext(
        ctx,
        &models,
        findAvailabilitiesByConsultantIDQuery,
        consultantID,
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