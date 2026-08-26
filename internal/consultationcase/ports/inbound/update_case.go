package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/usecase/dto"
)

type UpdateCaseInt interface {
	Execute(ctx context.Context, clientID string, caseID string, req *dto.UpdateCaseDTO) error
}