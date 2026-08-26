package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

var (
	ErrUnauthorizedAccess = custom_errors.UnauthorizedException("case does not belong to the client")
	ErrCaseNotFound = custom_errors.NotFoundError("case not found")
)