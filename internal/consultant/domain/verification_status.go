package domain

type VerificationStatus string

const (
	VerificationStatusNotStarted VerificationStatus = "NOT_STARTED"
	VerificationStatusPending  VerificationStatus = "PENDING"
	VerificationStatusReviewing VerificationStatus = "REVIEW"
	VerificationStatusApproved VerificationStatus = "APPROVED"
	VerificationStatusRejected VerificationStatus = "REJECTED"
)

func (v VerificationStatus) IsValid() bool {
	switch v {
	case VerificationStatusNotStarted, VerificationStatusPending, VerificationStatusReviewing, VerificationStatusApproved, VerificationStatusRejected:
		return true
	default:
		return false
	}
}

func NewVerificationStatus(status string) (VerificationStatus, error) {
	vs := VerificationStatus(status)
	if !vs.IsValid() {
		return "", ErrInvalidVerificationStatus
	}
	return vs, nil
}