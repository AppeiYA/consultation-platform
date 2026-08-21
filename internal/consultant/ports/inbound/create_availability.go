package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
)

type CreateAvailabilityInt interface {
	Execute(ctx context.Context, userID string, req *dto.CreateAvailabilityRequest) error
}