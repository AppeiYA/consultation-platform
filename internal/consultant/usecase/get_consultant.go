package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
)

type GetConsultant struct {
	consultantRepo outbound.ConsultantRepository
}

func NewGetConsultantUsecase(consultantRepo outbound.ConsultantRepository) *GetConsultant {
	return &GetConsultant{
		consultantRepo: consultantRepo,
	}
}

func(uc *GetConsultant) ByID(ctx context.Context, id string) (*domain.Consultant, error) {
	return uc.consultantRepo.FindByID(ctx, id)
}

func (uc *GetConsultant) ByUserID(ctx context.Context, userID string) (*domain.Consultant, error) {
	return uc.consultantRepo.FindByUserID(ctx, userID)
}