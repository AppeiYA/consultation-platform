package consultantAvailabilityRepo

import "context"

var deleteAvailabilityQuery = `
DELETE FROM consultant_availabilities WHERE id = $1
`
func (a *AvailabilityRepository) DeleteAvailability(ctx context.Context, availabilityID string) error {
	executor := a.repository.Executor(ctx)
	_, err := executor.ExecContext(ctx, deleteAvailabilityQuery, availabilityID)
	if err != nil {
		return err
	}
	return nil
}