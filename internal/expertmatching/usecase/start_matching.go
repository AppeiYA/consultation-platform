package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/inbound"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

type StartMatchingUsecase struct {
	caseReader    outbound.CaseReader
	runRepository outbound.MatchingRunRepository
	dispatcher    outbound.MatchingJobDispatcher
	idGenerator   shared_outbound.IdentifierGenerator
	clock         shared_outbound.Clock
}

func NewStartMatchingUsecase(
	caseReader outbound.CaseReader,
	runRepository outbound.MatchingRunRepository,
	dispatcher outbound.MatchingJobDispatcher,
	idGenerator shared_outbound.IdentifierGenerator,
	clock shared_outbound.Clock,
) *StartMatchingUsecase {
	return &StartMatchingUsecase{
		caseReader:    caseReader,
		runRepository: runRepository,
		dispatcher:    dispatcher,
		idGenerator:   idGenerator,
		clock:         clock,
	}
}

func (uc *StartMatchingUsecase) Execute(ctx context.Context, caseID string) (*domain.MatchingRun, error) {
	if len(caseID) == 0 {
		return nil, domain.ErrInvalidCaseID
	}

	// Verify case exists
	_, err := uc.caseReader.GetCaseDetails(ctx, caseID)
	if err != nil {
		return nil, err
	}

	rankingVersion, err := domain.NewRankingVersion("v1")
	if err != nil {
		return nil, err
	}

	runID, err := uc.idGenerator.Generate("mrun")
	if err != nil {
		return nil, err
	}

	now := uc.clock.Now()
	run, err := domain.NewMatchingRun(runID, caseID, rankingVersion, now)
	if err != nil {
		return nil, err
	}

	if err := uc.runRepository.Save(ctx, &run); err != nil {
		return nil, err
	}

	if err := uc.dispatcher.DispatchMatching(ctx, runID); err != nil {
		_ = run.Fail("failed to dispatch matching job: "+err.Error(), uc.clock.Now())
		_ = uc.runRepository.Save(ctx, &run)
		return nil, err
	}

	return &run, nil
}

var _ inbound.StartMatchingInt = (*StartMatchingUsecase)(nil)
