package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/inbound"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

type GetMatchingResultUsecase struct {
	runRepository outbound.MatchingRunRepository
}

func NewGetMatchingResultUsecase(runRepository outbound.MatchingRunRepository) *GetMatchingResultUsecase {
	return &GetMatchingResultUsecase{
		runRepository: runRepository,
	}
}

func (uc *GetMatchingResultUsecase) Execute(ctx context.Context, req inbound.GetMatchingResultRequest) (*inbound.MatchingResultResponse, error) {
	var run *domain.MatchingRun
	var err error

	if len(req.RunID) > 0 {
		run, err = uc.runRepository.FindByID(ctx, req.RunID)
	} else if len(req.CaseID) > 0 {
		run, err = uc.runRepository.FindLatestByCaseID(ctx, req.CaseID)
	} else {
		return nil, domain.ErrInvalidCaseID
	}

	if err != nil {
		return nil, err
	}
	if run == nil {
		return nil, custom_errors.NotFoundError("matching run not found")
	}

	topN := req.TopN
	if topN <= 0 {
		topN = 5
	}

	topCandidates := run.TopN(topN)

	return &inbound.MatchingResultResponse{
		RunID:          run.ID(),
		CaseID:         run.CaseID(),
		Status:         string(run.Status()),
		RankingVersion: run.RankingVersion().Value(),
		TopCandidates:  topCandidates,
		TotalRanked:    len(run.Candidates()),
		CompletedAt:    run.CompletedAt(),
	}, nil
}
