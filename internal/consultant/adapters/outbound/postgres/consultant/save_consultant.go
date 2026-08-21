package consultantRepo

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var saveConsultantQuery = `
	INSERT INTO consultants (
		id,
		user_id,
		profession_id,
		display_name,
		bio,
		years_experience,
		is_accepting_clients,
		created_at,
		updated_at
	)
	VALUES (
		$1, $2, $3, $4, $5,
		$6, $7, $8, $9
	)
`

func (r *ConsultantRepository) Save(ctx context.Context, consultant *domain.Consultant) error {
	var executor = r.repository.Executor(ctx)
	model := ConsultantFromDomain(consultant)

	_, err := executor.ExecContext(
		ctx, 
		saveConsultantQuery,
		model.ID,
		model.UserID,
		model.ProfessionID,
		model.DisplayName,
		model.Bio,
		model.YearsExperience,
		model.IsAcceptingClients,
		model.CreatedAt,
		model.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}