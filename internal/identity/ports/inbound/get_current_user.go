package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type GetCurrentUserInt interface {
	Execute(ctx context.Context, req dto.GetCurrentUserRequest) (dto.GetCurrentUserResponse, error)
}