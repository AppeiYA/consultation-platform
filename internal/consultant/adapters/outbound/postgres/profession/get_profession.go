package professionRepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var getProfessionByIDQuery = `
	SELECT id, name, created_at
	FROM professions
	WHERE id = $1
`
func (p *ProfessionRepository) GetProfessionByID(ctx context.Context, professionID string) (*domain.Profession, error) {
	executor := p.repository.Executor(ctx)
	var model ProfessionModel
	err := executor.GetContext(ctx, &model, getProfessionByIDQuery, professionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrInvalidProfession
		}
		return nil, err
	}

	domainProfession := model.ToDomain()

	return &domainProfession, nil
}