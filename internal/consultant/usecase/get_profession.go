package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
)

type GetProfessionUsecase struct {
	professionRepository outbound.ProfessionRepository
}

func NewGetProfessionUsecase(professionRepository outbound.ProfessionRepository) *GetProfessionUsecase {
	return &GetProfessionUsecase{
		professionRepository: professionRepository,
	}
}

func (g *GetProfessionUsecase) Execute(ctx context.Context, professionID string) (*domain.Profession, error) {
	profession, err := g.professionRepository.GetProfessionByID(ctx, professionID)
	if err != nil {
		return nil, err
	}

	return profession, nil
}