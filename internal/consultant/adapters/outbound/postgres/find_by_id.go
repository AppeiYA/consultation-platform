package postgres

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var findByIDQuery = `SELECT * FROM consultants WHERE id = $1`

func (r *ConsultantRepository) FindByID(ctx context.Context, id string) (*domain.Consultant, error) {
	var executor = r.repository.Executor(ctx)
	var model Consultant

	err := executor.GetContext(ctx, &model, findByIDQuery, id)
	if err != nil {
		return nil, err
	}

	return model.ToDomain()
}