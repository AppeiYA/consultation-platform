package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
)

type GetCaseInt interface {
	Execute(ctx context.Context, clientID string, caseID string) (*domain.ConsultationCase, error)
}