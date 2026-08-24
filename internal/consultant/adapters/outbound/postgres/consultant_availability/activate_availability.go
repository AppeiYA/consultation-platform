package consultantAvailabilityRepo

import "context"

var activateAvailabilitySQL = `
	UPDATE consultant_availabilities
	SET is_active = true, updated_at = NOW()
	WHERE id = $1;
`

func (r *AvailabilityRepository) ActivateAvailability(ctx context.Context, id string) error {
	executor := r.repository.Executor(ctx)

	_, err := executor.ExecContext(ctx, activateAvailabilitySQL, id)
	if err != nil {
		return err
	}

	return nil
}