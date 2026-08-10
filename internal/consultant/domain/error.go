package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

var (
	ErrConsultantAlreadyExists = custom_errors.BadException("consultant already exists")
)