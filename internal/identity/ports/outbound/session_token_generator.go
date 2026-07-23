package outbound

import "github.com/AppeiYA/consultation-platform/internal/identity/domain"

type SessionTokenGenerator interface {
    Generate() (domain.SessionToken, error)
}