package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

type UpdateConsultant struct {
	consultantRepo outbound.ConsultantRepository
	professionRepo outbound.ProfessionRepository
	clock          shared_outbound.Clock
}

func NewUpdateConsultantUsecase(
	consultantRepo outbound.ConsultantRepository,
	professionRepo outbound.ProfessionRepository,
	clock shared_outbound.Clock,
) *UpdateConsultant {
	return &UpdateConsultant{
		consultantRepo: consultantRepo,
		professionRepo: professionRepo,
		clock:          clock,
	}
}

func (uc *UpdateConsultant) Execute(ctx context.Context, userID string, input dto.UpdateConsultantDTO) error {
	consultant, err := uc.consultantRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if input.ProfessionID == "" {
		return domain.ErrInvalidProfession
	}

	profession, err := uc.professionRepo.GetProfessionByID(ctx, input.ProfessionID)
	if err != nil {
		return err
	}
	if profession == nil {
		return domain.ErrInvalidProfession
	}
	newDisplayName, err := domain.NewDisplayName(input.DisplayName)
	if err != nil {
		return err
	}
	newBio, err := domain.NewBio(input.Bio)
	if err != nil {
		return err
	}
	newYearsExperience, err := domain.NewYearsExperience(input.YearsExperience)
	if err != nil {
		return err
	}

	consultant.UpdateProfile(*profession, newDisplayName, newBio, newYearsExperience, uc.clock.Now())

	return uc.consultantRepo.Update(ctx, consultant)
}