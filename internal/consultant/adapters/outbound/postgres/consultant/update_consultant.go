package consultantRepo

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var updateConsultantQuery = `
	UPDATE consultants
	SET 
		bio = :bio, 
		profession_id = :profession_id, 
		display_name = :display_name,
		years_experience = :years_experience,
		updated_at = :updated_at
	WHERE id = :id
`

func (r *ConsultantRepository) Update(
	ctx context.Context,
	consultant *domain.Consultant,
) error {
	executor := r.repository.Executor(ctx)

	model := ConsultantFromDomain(consultant)
	model.UpdatedAt = r.clock.Now()

	_, err := executor.NamedExecContext(
		ctx,
		updateConsultantQuery,
		model,
	)

	return err
}