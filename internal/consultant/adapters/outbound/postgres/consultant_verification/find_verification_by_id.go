package consultantVerificationRepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var findVerificationByID = `
	SELECT
		id,
		consultant_id,
		provider,
		provider_reference,
		status,
		submitted_at,
		completed_at,
		created_at,
		updated_at
	FROM consultant_verifications
	WHERE id = $1
`
func (v *VerificationRepository) FindByID(ctx context.Context, verificationID string) (*domain.ConsultantVerification, error) {
	executor := v.repository.Executor(ctx)
	var model ConsultantVerification

	err := executor.GetContext(
		ctx, 
		&model,
		findVerificationByID,
		verificationID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrConsultantVerificationNotFound
		}
		return nil, err
	}

	return model.ToDomain()
}