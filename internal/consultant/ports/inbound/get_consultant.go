package inbound

import (
	"context"


	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
)

type GetConsultantInt interface {
	ByID(ctx context.Context, id string) (*dto.GetConsultantResponseDto, error)
	ByUserID(ctx context.Context, userID string) (*dto.GetConsultantResponseDto, error)
}