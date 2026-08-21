package domain

import "time"

var (
	VerificationIDPrefix = "ver"
)

type ConsultantVerification struct {
	id string
	consultantID string
	provider string
	provider_reference string
	status VerificationStatus
	submittedAt *time.Time
	completedAt *time.Time

	createdAt time.Time
	updatedAt time.Time
}

func NewVerification(
	id string,
	consultantID string,
	provider string,
	provider_reference string,
	status VerificationStatus,
	now time.Time,
) *ConsultantVerification {
	return &ConsultantVerification{
		id: id,
		consultantID: consultantID,
		provider: provider,
		provider_reference: provider_reference,
		status: status,
		submittedAt: &now,
		createdAt: now,
		updatedAt: now,
	}
}

func ReconstitueConsultantVerification(
	id string,
	consultantID string,
	provider string,
	provider_reference string,
	status VerificationStatus,
	submittedAt *time.Time,
	completedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) (*ConsultantVerification, error) {
	if !status.IsValid() {
		return nil, ErrInvalidVerificationStatus
	}

	return &ConsultantVerification{
		id: id,
		consultantID: consultantID,
		provider: provider,
		provider_reference: provider_reference,
		status: status,
		submittedAt: submittedAt,
		completedAt: completedAt,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

// getters
func (v *ConsultantVerification) ID() string {
	return v.id
}

func (v *ConsultantVerification) ConsultantID() string {
	return v.consultantID
}

func (v *ConsultantVerification) Provider() string {
	return v.provider
}

func (v *ConsultantVerification) ProviderReference() string {
	return v.provider_reference
}

func (v *ConsultantVerification) Status() VerificationStatus {
	return v.status
}

func (v *ConsultantVerification) SubmittedAt() *time.Time {
	return v.submittedAt
}

func (v *ConsultantVerification) CompletedAt() *time.Time {
	return v.completedAt
}

func (v *ConsultantVerification) CreatedAt() time.Time {
	return v.createdAt
}

func (v *ConsultantVerification) UpdatedAt() time.Time {
	return v.updatedAt
}