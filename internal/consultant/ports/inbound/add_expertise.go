package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
)

type AddExpertiseInt interface {
	Execute(ctx context.Context, userID string, req dto.AddExpertiseDTO) (*dto.ExpertiseResponseDTO, error)
}
