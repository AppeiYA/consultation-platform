package consultantRepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var findByIDQuery = `SELECT * FROM consultants WHERE id = $1`

func (r *ConsultantRepository) FindByID(ctx context.Context, id string) (*domain.Consultant, error) {
	var executor = r.repository.Executor(ctx)
	var model Consultant

	err := executor.GetContext(ctx, &model, findByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrConsultantNotFound
		}
		return nil, err
	}

	return model.ToDomain()
}