package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
)

type SubmitVerificationInt interface {
	Execute(ctx context.Context, userID string) (*dto.SubmitVerificationResponse, error)
}