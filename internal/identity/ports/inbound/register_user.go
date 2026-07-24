package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/usecase"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type RegisterUserInt interface {
	Execute(ctx context.Context, params usecase.RegisterUserParams) (*dto.RegisterUserResponse, error)
}