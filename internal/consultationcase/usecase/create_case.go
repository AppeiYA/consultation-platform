package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/usecase/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
	"go.uber.org/zap"
)

type CreateCaseUsecase struct {
	caseRepository  outbound.CaseRepository
	idGenerator     shared_outbound.IdentifierGenerator
	clientVerifier  outbound.ClientVerifier
	matchingStarter outbound.MatchingStarter
	clock           shared_outbound.Clock
}

func NewCreateCaseUsecase(
	caseRepository outbound.CaseRepository,
	idGenerator shared_outbound.IdentifierGenerator,
	clientVerifier outbound.ClientVerifier,
	matchingStarter outbound.MatchingStarter,
	clock shared_outbound.Clock,
) *CreateCaseUsecase {
	return &CreateCaseUsecase{
		caseRepository:  caseRepository,
		idGenerator:     idGenerator,
		clientVerifier:  clientVerifier,
		matchingStarter: matchingStarter,
		clock:           clock,
	}
}

func (uc *CreateCaseUsecase) Execute(ctx context.Context, clientID string, req *dto.CreateCaseDTO) error {
	// confirm client exists
	if err := uc.clientVerifier.VerifyClient(ctx, clientID); err != nil {
		return err
	}

	// generate a new case ID
	caseID, err := uc.idGenerator.Generate(domain.CaseIDPrefix)
	if err != nil {
		return err
	}

	// create a new consultation case
	now := uc.clock.Now()
	caseTitle, err := domain.NewCaseTitle(req.Title)
	if err != nil {
		return err
	}
	caseDescription, err := domain.NewCaseDescription(req.Description)
	if err != nil {
		return err
	}
	caseCategory, err := domain.NewCaseCategory(req.Category)
	if err != nil {
		return err
	}

	newCase, err := domain.NewConsultationCase(
		caseID,
		clientID,
		*caseTitle,
		*caseDescription,
		*caseCategory,
		now,
	)
	if err != nil {
		return err
	}

	// save the new case to the repository
	if err := uc.caseRepository.SaveCase(ctx, newCase); err != nil {
		return err
	}

	// auto-trigger matching if a matching starter is configured
	if uc.matchingStarter != nil {
		if err := uc.matchingStarter.StartMatching(ctx, caseID); err != nil {
			logger.Warn(
				"failed to auto-trigger matching on case creation",
				zap.Error(err),
				zap.String("case_id", caseID),
			)
		} else {
			logger.Info(
				"successfully auto-triggered matching on case creation",
				zap.String("case_id", caseID),
			)
		}
	}

	return nil
}

