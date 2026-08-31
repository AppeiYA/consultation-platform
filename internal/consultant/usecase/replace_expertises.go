package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

type ReplaceExpertises struct {
	consultantRepo outbound.ConsultantRepository
	expertiseRepo  outbound.ExpertiseRepository
	idGenerator    shared_outbound.IdentifierGenerator
}

func NewReplaceExpertisesUsecase(
	consultantRepo outbound.ConsultantRepository,
	expertiseRepo outbound.ExpertiseRepository,
	idGenerator shared_outbound.IdentifierGenerator,
) *ReplaceExpertises {
	return &ReplaceExpertises{
		consultantRepo: consultantRepo,
		expertiseRepo:  expertiseRepo,
		idGenerator:    idGenerator,
	}
}

func (uc *ReplaceExpertises) Execute(ctx context.Context, userID string, req dto.ReplaceExpertisesDTO) ([]dto.ExpertiseResponseDTO, error) {
	consultant, err := uc.consultantRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	entities := make([]*domain.Expertise, 0, len(req.Expertises))
	for _, name := range req.Expertises {
		expID, err := uc.idGenerator.Generate(domain.ExpertiseIDPrefix)
		if err != nil {
			return nil, err
		}
		exp, err := domain.NewExpertise(expID, consultant.ID(), name)
		if err != nil {
			return nil, err
		}
		entities = append(entities, exp)
	}

	if err := uc.expertiseRepo.ReplaceAll(ctx, consultant.ID(), entities); err != nil {
		return nil, err
	}

	results := make([]dto.ExpertiseResponseDTO, 0, len(entities))
	for _, exp := range entities {
		results = append(results, dto.ExpertiseResponseDTO{
			ID:           exp.ID(),
			ConsultantID: exp.ConsultantID(),
			Name:         exp.Name(),
		})
	}

	return results, nil
}
