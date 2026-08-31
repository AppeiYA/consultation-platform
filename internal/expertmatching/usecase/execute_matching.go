package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/inbound"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

type ExecuteMatchingUsecase struct {
	caseReader         outbound.CaseReader
	candidateGenerator outbound.CandidateGenerator
	candidateRanker    outbound.CandidateRanker
	runRepository      outbound.MatchingRunRepository
	clock              shared_outbound.Clock
}

func NewExecuteMatchingUsecase(
	caseReader outbound.CaseReader,
	candidateGenerator outbound.CandidateGenerator,
	candidateRanker outbound.CandidateRanker,
	runRepository outbound.MatchingRunRepository,
	clock shared_outbound.Clock,
) *ExecuteMatchingUsecase {
	return &ExecuteMatchingUsecase{
		caseReader:         caseReader,
		candidateGenerator: candidateGenerator,
		candidateRanker:    candidateRanker,
		runRepository:      runRepository,
		clock:              clock,
	}
}

func (uc *ExecuteMatchingUsecase) Execute(ctx context.Context, runID string) (*domain.MatchingRun, error) {
	if len(runID) == 0 {
		return nil, domain.ErrInvalidMatchingRunID
	}

	run, err := uc.runRepository.FindByID(ctx, runID)
	if err != nil {
		return nil, err
	}

	caseDetails, err := uc.caseReader.GetCaseDetails(ctx, run.CaseID())
	if err != nil {
		_ = run.Fail(err.Error(), uc.clock.Now())
		_ = uc.runRepository.Save(ctx, run)
		return nil, err
	}

	// 1. Generation phase
	if err := run.StartGeneration(); err != nil {
		return nil, err
	}

	candidatePool, err := uc.candidateGenerator.GenerateCandidates(ctx, domain.NewCandidateGenerationCriteria(
		caseDetails.Category,
		false,
	))
	if err != nil {
		_ = run.Fail(err.Error(), uc.clock.Now())
		_ = uc.runRepository.Save(ctx, run)
		return nil, err
	}

	// 2. Ranking phase
	if err := run.StartRanking(); err != nil {
		return nil, err
	}

	rankedResults, err := uc.candidateRanker.Rank(ctx, outbound.RankingRequest{
		CaseDetails:    *caseDetails,
		RankingVersion: run.RankingVersion(),
		CandidatePool:  candidatePool,
	})
	if err != nil {
		_ = run.Fail(err.Error(), uc.clock.Now())
		_ = uc.runRepository.Save(ctx, run)
		return nil, err
	}

	// 3. Complete phase
	if err := run.Complete(rankedResults, uc.clock.Now()); err != nil {
		_ = run.Fail(err.Error(), uc.clock.Now())
		_ = uc.runRepository.Save(ctx, run)
		return nil, err
	}

	if err := uc.runRepository.Save(ctx, run); err != nil {
		return nil, err
	}

	return run, nil
}

var _ inbound.ExecuteMatchingInt = (*ExecuteMatchingUsecase)(nil)
