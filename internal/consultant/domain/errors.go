package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

var ErrBioOverflow = custom_errors.BadException("Bio is too long")
var ErrEmptyBio = custom_errors.BadException("Bio is empty")
var ErrInvalidDisplayNameLength = custom_errors.BadException("Display name length must be > 6 and < 20")
var ErrInvalidDisplayName = custom_errors.BadException("special characters not allowed in display name")
var ErrInvalidYearsExperience = custom_errors.BadException("Years of experience must be > 0 and < 70")
var ErrInvalidProfession = custom_errors.BadException("profession is not in system")