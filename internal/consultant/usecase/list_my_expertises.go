package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
)

type ListMyExpertises struct {
	consultantRepo outbound.ConsultantRepository
	expertiseRepo  outbound.ExpertiseRepository
}

func NewListMyExpertisesUsecase(
	consultantRepo outbound.ConsultantRepository,
	expertiseRepo outbound.ExpertiseRepository,
) *ListMyExpertises {
	return &ListMyExpertises{
		consultantRepo: consultantRepo,
		expertiseRepo:  expertiseRepo,
	}
}

func (uc *ListMyExpertises) Execute(ctx context.Context, userID string) ([]dto.ExpertiseResponseDTO, error) {
	consultant, err := uc.consultantRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	expertises, err := uc.expertiseRepo.FindByConsultantID(ctx, consultant.ID())
	if err != nil {
		return nil, err
	}

	results := make([]dto.ExpertiseResponseDTO, 0, len(expertises))
	for _, exp := range expertises {
		results = append(results, dto.ExpertiseResponseDTO{
			ID:           exp.ID(),
			ConsultantID: exp.ConsultantID(),
			Name:         exp.Name(),
		})
	}

	return results, nil
}
