package consultationCaseRepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
)

var findCaseByIDQuery = `
SELECT id, client_id, title, description, category, status, created_at, updated_at
FROM consultation_cases
WHERE id = $1
`
func (r *ConsultationCaseRepository) FindCaseByID(ctx context.Context, caseID string) (*domain.ConsultationCase, error) {
	executor := r.repository.Executor(ctx)
	var model CaseModel
	err := executor.GetContext(ctx, &model, findCaseByIDQuery, caseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCaseNotFound
		}
		return nil, err
	}

	return model.ToDomain()
}