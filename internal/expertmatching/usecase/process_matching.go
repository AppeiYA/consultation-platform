package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/inbound"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
	"go.uber.org/zap"
)

type ProcessMatchingUsecase struct {
	caseReader         outbound.CaseReader
	candidateGenerator outbound.CandidateGenerator
	candidateRanker    outbound.CandidateRanker
	runRepository      outbound.MatchingRunRepository
	clock              shared_outbound.Clock
}

func NewProcessMatchingUsecase(
	caseReader outbound.CaseReader,
	candidateGenerator outbound.CandidateGenerator,
	candidateRanker outbound.CandidateRanker,
	runRepository outbound.MatchingRunRepository,
	clock shared_outbound.Clock,
) *ProcessMatchingUsecase {
	return &ProcessMatchingUsecase{
		caseReader:         caseReader,
		candidateGenerator: candidateGenerator,
		candidateRanker:    candidateRanker,
		runRepository:      runRepository,
		clock:              clock,
	}
}

func (uc *ProcessMatchingUsecase) Execute(ctx context.Context, runID string) error {
	if len(runID) == 0 {
		return domain.ErrInvalidMatchingRunID
	}

	run, err := uc.runRepository.FindByID(ctx, runID)
	if err != nil {
		return err
	}

	// Idempotency check: if run already completed or cancelled, do not re-process
	if run.IsCompleted() || run.Status() == domain.RunStatusCancelled {
		logger.Info(
			"matching run already completed or cancelled, skipping",
			zap.String("run_id", runID),
			zap.String("status", string(run.Status())),
		)
		return nil
	}

	caseDetails, err := uc.caseReader.GetCaseDetails(ctx, run.CaseID())
	if err != nil {
		logger.Error("failed to get case details for matching run", zap.Error(err), zap.String("run_id", runID), zap.String("case_id", run.CaseID()))
		_ = run.Fail(err.Error(), uc.clock.Now())
		_ = uc.runRepository.Save(ctx, run)
		return err
	}

	// 1. Generation phase
	if err := run.StartGeneration(); err != nil {
		return err
	}

	candidatePool, err := uc.candidateGenerator.GenerateCandidates(ctx, domain.NewCandidateGenerationCriteria(
		caseDetails.Category,
		false,
	))
	if err != nil {
		logger.Error("failed to generate candidates", zap.Error(err), zap.String("run_id", runID))
		_ = run.Fail(err.Error(), uc.clock.Now())
		_ = uc.runRepository.Save(ctx, run)
		return err
	}

	if len(candidatePool.Candidates()) == 0 {
		logger.Warn(
			"no matching candidates found for consultation case",
			zap.String("run_id", runID),
			zap.String("case_id", run.CaseID()),
			zap.String("category", caseDetails.Category.Value()),
		)
	} else {
		logger.Info(
			"candidates generated for matching run",
			zap.String("run_id", runID),
			zap.Int("candidate_pool_size", len(candidatePool.Candidates())),
		)
	}

	// 2. Ranking phase
	if err := run.StartRanking(); err != nil {
		return err
	}

	rankedResults, err := uc.candidateRanker.Rank(ctx, outbound.RankingRequest{
		CaseDetails:    *caseDetails,
		RankingVersion: run.RankingVersion(),
		CandidatePool:  candidatePool,
	})
	if err != nil {
		logger.Error("failed to rank candidates", zap.Error(err), zap.String("run_id", runID))
		_ = run.Fail(err.Error(), uc.clock.Now())
		_ = uc.runRepository.Save(ctx, run)
		return err
	}

	logger.Info(
		"candidate ranking completed",
		zap.String("run_id", runID),
		zap.Int("ranked_candidates_count", len(rankedResults)),
	)

	// 3. Complete phase
	if err := run.Complete(rankedResults, uc.clock.Now()); err != nil {
		logger.Error("failed to complete matching run", zap.Error(err), zap.String("run_id", runID))
		_ = run.Fail(err.Error(), uc.clock.Now())
		_ = uc.runRepository.Save(ctx, run)
		return err
	}

	if err := uc.runRepository.Save(ctx, run); err != nil {
		logger.Error("failed to save completed matching run", zap.Error(err), zap.String("run_id", runID))
		return err
	}

	logger.Info(
		"matching run completed successfully",
		zap.String("run_id", runID),
		zap.String("case_id", run.CaseID()),
		zap.Int("total_ranked", len(rankedResults)),
	)

	return nil
}

var _ inbound.ProcessMatchingInt = (*ProcessMatchingUsecase)(nil)
