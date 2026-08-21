package consultantRepo

import "context"

var existsByUserIDQuery = `SELECT EXISTS(SELECT 1 FROM consultants WHERE user_id = $1)`

func (r *ConsultantRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	var exists bool
	var executor = r.repository.Executor(ctx)
	err := executor.QueryRowxContext(ctx, existsByUserIDQuery, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}