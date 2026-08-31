package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
)

type GetConsultant struct {
	consultantRepo outbound.ConsultantRepository
	professionRepo outbound.ProfessionRepository
	expertiseRepo  outbound.ExpertiseRepository
}

func NewGetConsultantUsecase(
	consultantRepo outbound.ConsultantRepository,
	professionRepo outbound.ProfessionRepository,
	expertiseRepo outbound.ExpertiseRepository,
) *GetConsultant {
	return &GetConsultant{
		consultantRepo: consultantRepo,
		professionRepo: professionRepo,
		expertiseRepo:  expertiseRepo,
	}
}

func (uc *GetConsultant) ByID(ctx context.Context, id string) (*dto.GetConsultantResponseDto, error) {
	consultant, err := uc.consultantRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	profession, err := uc.professionRepo.GetProfessionByID(ctx, consultant.ProfessionID())
	if err != nil {
		return nil, err
	}
	if profession == nil {
		return nil, domain.ErrInvalidProfession
	}

	expertises, err := uc.expertiseRepo.FindByConsultantID(ctx, consultant.ID())
	if err != nil {
		return nil, err
	}

	expertiseNames := make([]string, 0, len(expertises))
	for _, exp := range expertises {
		expertiseNames = append(expertiseNames, exp.Name())
	}

	return &dto.GetConsultantResponseDto{
		ID:                 consultant.ID(),
		Profession:         profession.Name(),
		DisplayName:        consultant.DisplayName().String(),
		UserID:             consultant.UserID(),
		Bio:                consultant.Bio().String(),
		YearsExperience:    consultant.YearsExperience().Int(),
		IsAcceptingClients: consultant.IsAcceptingClients(),
		Expertises:         expertiseNames,
		CreatedAt:          consultant.CreatedAt().Format("2006-01-02 15:04:05"),
		UpdatedAt:          consultant.UpdatedAt().Format("2006-01-02 15:04:05"),
	}, nil
}

func (uc *GetConsultant) ByUserID(ctx context.Context, userID string) (*dto.GetConsultantResponseDto, error) {
	consultant, err := uc.consultantRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	profession, err := uc.professionRepo.GetProfessionByID(ctx, consultant.ProfessionID())
	if err != nil {
		return nil, err
	}
	if profession == nil {
		return nil, domain.ErrInvalidProfession
	}

	expertises, err := uc.expertiseRepo.FindByConsultantID(ctx, consultant.ID())
	if err != nil {
		return nil, err
	}

	expertiseNames := make([]string, 0, len(expertises))
	for _, exp := range expertises {
		expertiseNames = append(expertiseNames, exp.Name())
	}

	return &dto.GetConsultantResponseDto{
		ID:                 consultant.ID(),
		Profession:         profession.Name(),
		DisplayName:        consultant.DisplayName().String(),
		UserID:             consultant.UserID(),
		Bio:                consultant.Bio().String(),
		YearsExperience:    consultant.YearsExperience().Int(),
		IsAcceptingClients: consultant.IsAcceptingClients(),
		Expertises:         expertiseNames,
		CreatedAt:          consultant.CreatedAt().Format("2006-01-02 15:04:05"),
		UpdatedAt:          consultant.UpdatedAt().Format("2006-01-02 15:04:05"),
	}, nil
}
