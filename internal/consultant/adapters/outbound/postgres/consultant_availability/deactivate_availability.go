package consultantAvailabilityRepo

import "context"

var deactivateAvailabilitySQL = `
	UPDATE consultant_availabilities
	SET is_active = false, updated_at = NOW()
	WHERE id = $1;
`

func (r *AvailabilityRepository) DeactivateAvailability(ctx context.Context, id string) error {
	executor := r.repository.Executor(ctx)

	_, err := executor.ExecContext(ctx, deactivateAvailabilitySQL, id)
	if err != nil {
		return err
	}

	return nil
}