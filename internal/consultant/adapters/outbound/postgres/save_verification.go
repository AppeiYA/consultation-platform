package postgres

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

var saveVerificationQuery = `
	INSERT INTO consultant_verifications (id, consultant_id, provider, provider_reference, status, submitted_at, completed_at, created_at, updated_at)
	VALUES (:id, :consultant_id, :provider, :provider_reference, :status, :submitted_at, :completed_at, :created_at, :updated_at)
`

func (v *VerificationRepository) Save(ctx context.Context, verification *domain.ConsultantVerification) error {
	executor := v.repository.Executor(ctx)
	model := ConsultantVerificationFromDomain(verification)

	_, err := executor.NamedExecContext(ctx, saveVerificationQuery, model)
	if err != nil {
		return err
	}

	return nil
}