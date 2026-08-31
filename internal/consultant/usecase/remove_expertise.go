package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
)

type RemoveExpertise struct {
	consultantRepo outbound.ConsultantRepository
	expertiseRepo  outbound.ExpertiseRepository
}

func NewRemoveExpertiseUsecase(
	consultantRepo outbound.ConsultantRepository,
	expertiseRepo outbound.ExpertiseRepository,
) *RemoveExpertise {
	return &RemoveExpertise{
		consultantRepo: consultantRepo,
		expertiseRepo:  expertiseRepo,
	}
}

func (uc *RemoveExpertise) Execute(ctx context.Context, userID string, expertiseID string) error {
	consultant, err := uc.consultantRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return uc.expertiseRepo.Delete(ctx, consultant.ID(), expertiseID)
}
