package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

type AddExpertise struct {
	consultantRepo outbound.ConsultantRepository
	expertiseRepo  outbound.ExpertiseRepository
	idGenerator    shared_outbound.IdentifierGenerator
}

func NewAddExpertiseUsecase(
	consultantRepo outbound.ConsultantRepository,
	expertiseRepo outbound.ExpertiseRepository,
	idGenerator shared_outbound.IdentifierGenerator,
) *AddExpertise {
	return &AddExpertise{
		consultantRepo: consultantRepo,
		expertiseRepo:  expertiseRepo,
		idGenerator:    idGenerator,
	}
}

func (uc *AddExpertise) Execute(ctx context.Context, userID string, req dto.AddExpertiseDTO) (*dto.ExpertiseResponseDTO, error) {
	consultant, err := uc.consultantRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	expID, err := uc.idGenerator.Generate(domain.ExpertiseIDPrefix)
	if err != nil {
		return nil, err
	}

	expertise, err := domain.NewExpertise(expID, consultant.ID(), req.Name)
	if err != nil {
		return nil, err
	}

	if err := uc.expertiseRepo.Add(ctx, expertise); err != nil {
		return nil, err
	}

	return &dto.ExpertiseResponseDTO{
		ID:           expertise.ID(),
		ConsultantID: expertise.ConsultantID(),
		Name:         expertise.Name(),
	}, nil
}
