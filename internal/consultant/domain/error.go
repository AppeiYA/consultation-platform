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

	ErrBioOverflow = custom_errors.BadException("Bio is too long")
	ErrEmptyBio = custom_errors.BadException("Bio is empty")

	ErrInvalidDisplayNameLength = custom_errors.BadException("Display name length must be > 6 and < 20")
	ErrInvalidDisplayName = custom_errors.BadException("special characters not allowed in display name")

	ErrInvalidYearsExperience = custom_errors.BadException("Years of experience must be > 0 and < 70")
	ErrInvalidProfession = custom_errors.BadException("profession is not in system")

	ErrConsultantNotAcceptingClients = custom_errors.BadException("consultant is not accepting clients")
)