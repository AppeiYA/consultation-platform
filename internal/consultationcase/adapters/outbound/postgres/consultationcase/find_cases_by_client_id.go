package consultationCaseRepo

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
)

var findCasesByClientIDQuery = `
SELECT id, client_id, title, description, category, status, created_at, updated_at
FROM consultation_cases
WHERE client_id = $1
`
func (r *ConsultationCaseRepository) FindCasesByClientID(ctx context.Context, clientID string) ([]*domain.ConsultationCase, error) {
	executor := r.repository.Executor(ctx)
	var models []CaseModel
	err := executor.SelectContext(ctx, &models, findCasesByClientIDQuery, clientID)
	if err != nil {
		return nil, err
	}

	cases := make([]*domain.ConsultationCase, len(models))
	for i, model := range models {
		c, err := model.ToDomain()
		if err != nil {
			return nil, err
		}
		cases[i] = c
	}

	return cases, nil
}