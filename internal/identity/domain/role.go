package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

var invalidRole = custom_errors.BadException("invalid role")

type Role string

const (
	RoleClient     Role = "CLIENT"
	RoleConsultant Role = "CONSULTANT"
	RoleAdmin      Role = "ADMIN"
)

func NewRole(role string) (Role, error) {
	r := Role(role)

	switch r {
	case RoleClient, RoleConsultant, RoleAdmin:
		return r, nil
	default:
		return "", invalidRole
	}
}

func (r Role) String() string {
	return string(r)
}