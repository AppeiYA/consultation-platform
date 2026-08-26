package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/ports/outbound"
)

type DeleteCaseUsecase struct {
	caseRepository outbound.CaseRepository
}

func NewDeleteCaseUsecase(caseRepository outbound.CaseRepository) *DeleteCaseUsecase {
	return &DeleteCaseUsecase{
		caseRepository: caseRepository,
	}
}

func (uc *DeleteCaseUsecase) Execute(ctx context.Context, clientID string, caseID string) error {
	// retrieve the consultation case
	consultationCase, err := uc.caseRepository.FindCaseByID(ctx, caseID)
	if err != nil {
		return err
	}

	if consultationCase == nil {
		return domain.ErrCaseNotFound
	}

	// confirm the case belongs to the client
	if consultationCase.ClientID() != clientID {
		return domain.ErrUnauthorizedAccess
	}

	// delete the consultation case
	if err := uc.caseRepository.DeleteCase(ctx, caseID); err != nil {
		return err
	}

	return nil
}