package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type LoginUserInt interface {
	Execute(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
}