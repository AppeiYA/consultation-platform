package outbound

import "context"

type ClientVerifier interface {
	VerifyClient(ctx context.Context, clientID string) error
}