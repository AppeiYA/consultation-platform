package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/usecase/dto"
)

type SaveCaseInt interface {
	Execute(ctx context.Context, clientID string, req *dto.CreateCaseDTO) error
}