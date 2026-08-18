package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

var (
	ErrConsultantAlreadyExists = custom_errors.ConflictError("consultant already exists")
	ErrConsultantNotFound      = custom_errors.NotFoundError("consultant not found")
	ErrInvalidVerificationStatus = custom_errors.BadException("invalid verification status")
	ErrVerificationPending = custom_errors.BadException("verification is pending")
	ErrVerificationInReview = custom_errors.BadException("verification is in review")
	ErrVerificationAlreadyApproved = custom_errors.BadException("verification is already approved")
	ErrConsultantVerificationNotFound = custom_errors.NotFoundError("consultant verification not found")
	ErrVerificationUnavailable = custom_errors.BadException("verification service unavailable")
)