package postgres

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var findByUserIDQuery = `SELECT * FROM consultants WHERE user_id = $1`

func (r *ConsultantRepository) FindByUserID(ctx context.Context, userID string) (*domain.Consultant, error) {
	var executor = r.repository.Executor(ctx)
	var model Consultant

	err := executor.GetContext(ctx, &model, findByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}

	return model.ToDomain()
}