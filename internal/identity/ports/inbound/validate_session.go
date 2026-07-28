package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type ValidateSessionInt interface {
	Execute(ctx context.Context, req dto.ValidateSessionRequest) (dto.ValidateSessionResponse, error)
}