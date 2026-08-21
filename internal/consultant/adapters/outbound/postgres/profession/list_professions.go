package professionRepo

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var getAllProfessionsQuery = `
	SELECT id, name, created_at
	FROM professions
`
func (r *ProfessionRepository) GetAllProfessions(ctx context.Context) ([]*domain.Profession, error) {
	executor := r.repository.Executor(ctx)

	var model []ProfessionModel
	err := executor.SelectContext(ctx, &model, getAllProfessionsQuery)
	if err != nil {
		return nil, err
	}

	var domainProfessions []*domain.Profession
	for _, m := range model {
		domainProfession := m.ToDomain()
		domainProfessions = append(domainProfessions, &domainProfession)
	}

	return domainProfessions, nil
}