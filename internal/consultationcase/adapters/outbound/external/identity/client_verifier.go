package consultationCaseIdentityAdapter

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type ClientVerifier struct {
	identityModule *identity.Module
}

func NewClientVerifier(
	identityModule *identity.Module,
) *ClientVerifier {
	return &ClientVerifier{
		identityModule: identityModule,
	}
}

func (c *ClientVerifier) VerifyClient(ctx context.Context, clientID string) error {
	_, err := c.identityModule.GetCurrentUser.Execute(ctx, dto.GetCurrentUserRequest{
		UserID: clientID,
	})
	if err != nil {
		return err
	}

	return nil
}