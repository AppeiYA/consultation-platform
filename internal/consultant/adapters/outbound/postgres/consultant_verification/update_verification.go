package consultantVerificationRepo

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var updateVerification = `
	UPDATE consultant_verifications
	SET
		provider = :provider,
		provider_reference = :provider_reference,
		status = :status,
		submitted_at = :submitted_at,
		completed_at = :completed_at,
		updated_at = :updated_at
	WHERE id = :id
`
func (v *VerificationRepository) Update(ctx context.Context, verification *domain.ConsultantVerification) error {
	executor := v.repository.Executor(ctx)
	model := ConsultantVerificationFromDomain(verification)

	_, err := executor.NamedExecContext(ctx, updateVerification, model)
	if err != nil {
		return err
	}

	return nil
}