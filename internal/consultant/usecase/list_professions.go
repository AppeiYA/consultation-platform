package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
)

type GetProfessionsUsecase struct {
	professionRepository outbound.ProfessionRepository
}

func NewListProfessionsUsecase(professionRepository outbound.ProfessionRepository) *GetProfessionsUsecase {
	return &GetProfessionsUsecase{
		professionRepository: professionRepository,
	}
}

func (g *GetProfessionsUsecase) Execute(ctx context.Context) ([]*domain.Profession, error) {
	professions, err := g.professionRepository.GetAllProfessions(ctx)
	if err != nil {
		return nil, err
	}

	return professions, nil
}