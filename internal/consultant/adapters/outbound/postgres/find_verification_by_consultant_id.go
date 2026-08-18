package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var findVerificationByConsultantIDQuery = `
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
	WHERE consultant_id = $1
	ORDER BY created_at DESC
	LIMIT 1
`
func (v *VerificationRepository) FindByConsultantID(
	ctx context.Context,
	consultantID string,
) (*domain.ConsultantVerification, error) {
	executor := v.repository.Executor(ctx)

	var model ConsultantVerification

	err := executor.GetContext(
		ctx,
		&model,
		findVerificationByConsultantIDQuery,
		consultantID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrConsultantVerificationNotFound
		}
		return nil, err
	}
	
	return model.ToDomain()
}