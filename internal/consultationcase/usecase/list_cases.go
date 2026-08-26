package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/ports/outbound"
)

type ListCasesUsecase struct {
	caseRepository outbound.CaseRepository
}

func NewListCasesUsecase(caseRepository outbound.CaseRepository) *ListCasesUsecase {
	return &ListCasesUsecase{
		caseRepository: caseRepository,
	}
}

func (uc *ListCasesUsecase) Execute(ctx context.Context, clientID string) ([]*domain.ConsultationCase, error) {
	// retrieve the consultation cases for the client
	cases, err := uc.caseRepository.FindCasesByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return cases, nil
}