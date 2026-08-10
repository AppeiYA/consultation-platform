package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
)

type RegisterConsultantInt interface {
	Execute(ctx context.Context, userID string, req *dto.RegisterConsultantDTO) error
}