package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
)

type UpdateConsultantInt interface {
	Execute(ctx context.Context, userID string, input dto.UpdateConsultantDTO) error
}