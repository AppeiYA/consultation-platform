package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type RegisterUserInt interface {
	Execute(ctx context.Context, params dto.RegisterUserRequest) (*dto.RegisterUserResponse, error)
}