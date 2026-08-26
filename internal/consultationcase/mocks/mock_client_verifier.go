package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/ports/outbound"
)

type MockClientVerifier struct {
	VerifyClientFn func(ctx context.Context, clientID string) error
}

func (m *MockClientVerifier) VerifyClient(ctx context.Context, clientID string) error {
	if m.VerifyClientFn != nil {
		return m.VerifyClientFn(ctx, clientID)
	}
	return nil
}

var _ outbound.ClientVerifier = (*MockClientVerifier)(nil)
