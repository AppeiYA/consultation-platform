package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/usecase/dto"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

type UpdateCaseUsecase struct {
	caseRepository outbound.CaseRepository
	clock shared_outbound.Clock
}

func NewUpdateCaseUsecase(
	caseRepository outbound.CaseRepository,
	clock shared_outbound.Clock,
) *UpdateCaseUsecase {
	return &UpdateCaseUsecase{
		caseRepository: caseRepository,
		clock: clock,
	}
}

func (uc *UpdateCaseUsecase) Execute(ctx context.Context, clientID string, caseID string, req *dto.UpdateCaseDTO) error {
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

	if req.Title != nil {
		caseTitle, err := domain.NewCaseTitle(*req.Title)
		if err != nil {
			return err
		}
		consultationCase.UpdateTitle(*caseTitle)
	}

	if req.Description != nil {
		caseDescription, err := domain.NewCaseDescription(*req.Description)
		if err != nil {
			return err
		}
		consultationCase.UpdateDescription(*caseDescription)
	}

	if req.Category != nil {
		caseCategory, err := domain.NewCaseCategory(*req.Category)
		if err != nil {
			return err
		}
		consultationCase.UpdateCategory(*caseCategory)
	}
	consultationCase.EndUpdating(uc.clock.Now())

	if err := uc.caseRepository.UpdateCase(ctx, consultationCase); err != nil {
		return err
	}

	return nil
}