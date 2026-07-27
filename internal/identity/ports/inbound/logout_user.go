package inbound 

import (
	"context"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type LogoutUserInt interface {
	Execute(ctx context.Context, req dto.LogoutRequest) (dto.LogoutResponse, error)
}