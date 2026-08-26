package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
)

type ListCaseInt interface {
	Execute(ctx context.Context, clientID string) ([]*domain.ConsultationCase, error)
}