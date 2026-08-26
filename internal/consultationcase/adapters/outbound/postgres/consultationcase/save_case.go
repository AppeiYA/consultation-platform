package consultationCaseRepo

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
)

var saveCaseQuery = `INSERT INTO consultation_cases (id, client_id, title, description, category, status, created_at, updated_at)
	VALUES (:id, :client_id, :title, :description, :category, :status, :created_at, :updated_at)
	ON CONFLICT (id) DO UPDATE SET
		client_id = EXCLUDED.client_id,
		title = EXCLUDED.title,
		description = EXCLUDED.description,
		category = EXCLUDED.category,
		status = EXCLUDED.status,
		updated_at = EXCLUDED.updated_at
`
func (r *ConsultationCaseRepository) SaveCase(ctx context.Context, consultationCase *domain.ConsultationCase) error {
	executor := r.repository.Executor(ctx)
	model := FromDomainToModel(consultationCase)
	_, err := executor.NamedExecContext(ctx, saveCaseQuery, model)
	if err != nil {
		return err
	}
	return nil
}