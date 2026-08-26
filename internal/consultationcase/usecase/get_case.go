package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/ports/outbound"
)

type GetCaseUsecase struct {
	caseRepository outbound.CaseRepository
}

func NewGetCaseUsecase(caseRepository outbound.CaseRepository) *GetCaseUsecase {
	return &GetCaseUsecase{
		caseRepository: caseRepository,
	}
}

func (uc *GetCaseUsecase) Execute(ctx context.Context, clientID string, caseID string) (*domain.ConsultationCase, error) {
	// retrieve the consultation case
	consultationCase, err := uc.caseRepository.FindCaseByID(ctx, caseID)
	if err != nil {
		return nil, err
	}


	if consultationCase == nil {
		return nil, domain.ErrCaseNotFound
	}

	// confirm the case belongs to the client
	if consultationCase.ClientID() != clientID {
		return nil, domain.ErrUnauthorizedAccess
	}

	return consultationCase, nil
}