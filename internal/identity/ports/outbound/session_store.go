package outbound

import (
	"context"
	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
)

type SessionStore interface {
    Save(ctx context.Context, session *domain.Session) error

    FindByTokenHash(
        ctx context.Context,
        tokenHash string,
    ) (*domain.Session, error)

    Delete(ctx context.Context, tokenHash string) error
}