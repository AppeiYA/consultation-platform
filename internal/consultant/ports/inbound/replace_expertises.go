package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
)

type ReplaceExpertisesInt interface {
	Execute(ctx context.Context, userID string, req dto.ReplaceExpertisesDTO) ([]dto.ExpertiseResponseDTO, error)
}
