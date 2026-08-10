package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var findByUserIDQuery = `SELECT * FROM consultants WHERE user_id = $1`

func (r *ConsultantRepository) FindByUserID(ctx context.Context, userID string) (*domain.Consultant, error) {
	var executor = r.repository.Executor(ctx)
	var model Consultant

	err := executor.GetContext(ctx, &model, findByUserIDQuery, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrConsultantNotFound
		}
		return nil, err
	}

	return model.ToDomain()
}