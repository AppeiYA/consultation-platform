package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
)

type ListMyExpertisesInt interface {
	Execute(ctx context.Context, userID string) ([]dto.ExpertiseResponseDTO, error)
}
