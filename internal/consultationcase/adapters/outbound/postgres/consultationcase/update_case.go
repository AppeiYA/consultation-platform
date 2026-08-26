package consultationCaseRepo

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
)

var updateCaseQuery = `
UPDATE consultation_cases
SET title = :title, description = :description, category = :category, updated_at = :updated_at
WHERE id = :id
`
func (r *ConsultationCaseRepository) UpdateCase(ctx context.Context, consultationCase *domain.ConsultationCase) error {
	executor := r.repository.Executor(ctx)
	model := FromDomainToModel(consultationCase)
	_, err := executor.NamedExecContext(ctx, updateCaseQuery, model)
	if err != nil {
		return err
	}
	return nil
}