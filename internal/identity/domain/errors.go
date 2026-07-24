package domain

import (
	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

var (
	ErrUserNotFound    = custom_errors.NotFoundError("user not found")
	ErrUserAlreadyExists = custom_errors.ConflictError("user already exists")
	ErrInvalidPassword = custom_errors.UnauthorizedException("invalid password")
)