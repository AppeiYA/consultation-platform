package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/mocks"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
	shared_mocks "github.com/AppeiYA/consultation-platform/internal/shared/mocks"
	"github.com/stretchr/testify/require"
)

func TestProcessMatchingUsecase(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 8, 26, 12, 0, 0, 0, time.UTC)
	cat, _ := domain.NewMatchingCategory("Software Engineering")
	rankingVer, _ := domain.NewRankingVersion("v1")

	setup := func() (
		*mocks.MockCaseReader,
		*mocks.MockCandidateGenerator,
		*mocks.MockCandidateRanker,
		*mocks.MockMatchingRunRepository,
		*shared_mocks.MockClock,
		*ProcessMatchingUsecase,
	) {
		caseReader := &mocks.MockCaseReader{}
		candidateGen := &mocks.MockCandidateGenerator{}
		candidateRanker := &mocks.MockCandidateRanker{}
		runRepo := &mocks.MockMatchingRunRepository{}
		clock := &shared_mocks.MockClock{
			NowFn: func() time.Time { return now },
		}

		uc := NewProcessMatchingUsecase(
			caseReader,
			candidateGen,
			candidateRanker,
			runRepo,
			clock,
		)
		return caseReader, candidateGen, candidateRanker, runRepo, clock, uc
	}

	t.Run("successful end-to-end background matching execution", func(t *testing.T) {
		caseReader, candidateGen, candidateRanker, runRepo, _, uc := setup()

		run, _ := domain.NewMatchingRun("mrun_001", "case_001", rankingVer, now)
		runRepo.FindByIDFn = func(ctx context.Context, id string) (*domain.MatchingRun, error) {
			require.Equal(t, "mrun_001", id)
			return &run, nil
		}

		caseReader.GetCaseDetailsFn = func(ctx context.Context, caseID string) (*outbound.CaseDetails, error) {
			return &outbound.CaseDetails{
				CaseID:      caseID,
				ClientID:    "client_1",
				Category:    cat,
				Title:       "Backend Architecture Review",
				Description: "Need review for microservices",
			}, nil
		}

		candidateGen.GenerateCandidatesFn = func(ctx context.Context, criteria domain.CandidateGenerationCriteria) (domain.CandidatePool, error) {
			require.Equal(t, "Software Engineering", criteria.Category().Value())
			p1, _ := domain.NewCandidateProfile("con_1", cat, "Software Engineering", nil, 10, "bio 1")
			p2, _ := domain.NewCandidateProfile("con_2", cat, "Software Engineering", nil, 8, "bio 2")
			return domain.NewCandidatePool([]domain.CandidateProfile{p1, p2}, 100)
		}

		candidateRanker.RankFn = func(ctx context.Context, req outbound.RankingRequest) ([]domain.RankedCandidate, error) {
			require.Equal(t, 2, req.CandidatePool.Size())
			r1, _ := domain.NewRank(1)
			s1, _ := domain.NewMatchScore(0.95)
			r2, _ := domain.NewRank(2)
			s2, _ := domain.NewMatchScore(0.85)

			c1, _ := domain.NewRankedCandidate("con_1", r1, s1, nil)
			c2, _ := domain.NewRankedCandidate("con_2", r2, s2, nil)
			return []domain.RankedCandidate{c1, c2}, nil
		}

		var savedRun *domain.MatchingRun
		runRepo.SaveFn = func(ctx context.Context, r *domain.MatchingRun) error {
			savedRun = r
			return nil
		}

		err := uc.Execute(ctx, "mrun_001")

		require.NoError(t, err)
		require.NotNil(t, savedRun)
		require.Equal(t, "mrun_001", savedRun.ID())
		require.Equal(t, domain.RunStatusCompleted, savedRun.Status())
		require.True(t, savedRun.IsCompleted())
		require.Len(t, savedRun.Candidates(), 2)
	})

	t.Run("idempotency: no-op when run is already completed", func(t *testing.T) {
		caseReader, candidateGen, candidateRanker, runRepo, _, uc := setup()

		run, _ := domain.NewMatchingRun("mrun_done", "case_001", rankingVer, now)
		_ = run.StartGeneration()
		_ = run.StartRanking()
		_ = run.Complete([]domain.RankedCandidate{}, now)

		runRepo.FindByIDFn = func(ctx context.Context, id string) (*domain.MatchingRun, error) {
			return &run, nil
		}

		// These should NOT be called
		caseReaderCalled := false
		caseReader.GetCaseDetailsFn = func(ctx context.Context, caseID string) (*outbound.CaseDetails, error) {
			caseReaderCalled = true
			return nil, nil
		}
		candidateGenCalled := false
		candidateGen.GenerateCandidatesFn = func(ctx context.Context, criteria domain.CandidateGenerationCriteria) (domain.CandidatePool, error) {
			candidateGenCalled = true
			return domain.CandidatePool{}, nil
		}

		err := uc.Execute(ctx, "mrun_done")
		require.NoError(t, err)
		require.False(t, caseReaderCalled)
		require.False(t, candidateGenCalled)
		require.Nil(t, candidateRanker.RankFn)
	})

	t.Run("should reject empty run ID", func(t *testing.T) {
		_, _, _, _, _, uc := setup()
		err := uc.Execute(ctx, "")
		require.Equal(t, domain.ErrInvalidMatchingRunID, err)
	})

	t.Run("should fail when run repository fails to find run", func(t *testing.T) {
		_, _, _, runRepo, _, uc := setup()
		runRepo.FindByIDFn = func(ctx context.Context, id string) (*domain.MatchingRun, error) {
			return nil, errors.New("run not found")
		}

		err := uc.Execute(ctx, "mrun_404")
		require.Error(t, err)
		require.Equal(t, "run not found", err.Error())
	})

	t.Run("should record failure and save when candidate generator fails", func(t *testing.T) {
		caseReader, candidateGen, _, runRepo, _, uc := setup()

		run, _ := domain.NewMatchingRun("mrun_fail_gen", "case_001", rankingVer, now)
		runRepo.FindByIDFn = func(ctx context.Context, id string) (*domain.MatchingRun, error) {
			return &run, nil
		}
		caseReader.GetCaseDetailsFn = func(ctx context.Context, caseID string) (*outbound.CaseDetails, error) {
			return &outbound.CaseDetails{CaseID: caseID, Category: cat}, nil
		}
		candidateGen.GenerateCandidatesFn = func(ctx context.Context, criteria domain.CandidateGenerationCriteria) (domain.CandidatePool, error) {
			return domain.CandidatePool{}, errors.New("generator database error")
		}

		var savedRun *domain.MatchingRun
		runRepo.SaveFn = func(ctx context.Context, r *domain.MatchingRun) error {
			savedRun = r
			return nil
		}

		err := uc.Execute(ctx, "mrun_fail_gen")
		require.Error(t, err)
		require.Equal(t, "generator database error", err.Error())
		require.NotNil(t, savedRun)
		require.Equal(t, domain.RunStatusFailed, savedRun.Status())
		require.Equal(t, "generator database error", savedRun.FailureReason())
	})

	t.Run("should record failure and save when candidate ranker fails", func(t *testing.T) {
		caseReader, candidateGen, candidateRanker, runRepo, _, uc := setup()

		run, _ := domain.NewMatchingRun("mrun_fail_rank", "case_001", rankingVer, now)
		runRepo.FindByIDFn = func(ctx context.Context, id string) (*domain.MatchingRun, error) {
			return &run, nil
		}
		caseReader.GetCaseDetailsFn = func(ctx context.Context, caseID string) (*outbound.CaseDetails, error) {
			return &outbound.CaseDetails{CaseID: caseID, Category: cat}, nil
		}
		candidateGen.GenerateCandidatesFn = func(ctx context.Context, criteria domain.CandidateGenerationCriteria) (domain.CandidatePool, error) {
			p, _ := domain.NewCandidateProfile("con_1", cat, "Tech", nil, 5, "bio")
			return domain.NewCandidatePool([]domain.CandidateProfile{p}, 100)
		}
		candidateRanker.RankFn = func(ctx context.Context, req outbound.RankingRequest) ([]domain.RankedCandidate, error) {
			return nil, errors.New("ranking model timeout")
		}

		var savedRun *domain.MatchingRun
		runRepo.SaveFn = func(ctx context.Context, r *domain.MatchingRun) error {
			savedRun = r
			return nil
		}

		err := uc.Execute(ctx, "mrun_fail_rank")
		require.Error(t, err)
		require.Equal(t, "ranking model timeout", err.Error())
		require.NotNil(t, savedRun)
		require.Equal(t, domain.RunStatusFailed, savedRun.Status())
		require.Equal(t, "ranking model timeout", savedRun.FailureReason())
	})

	t.Run("completes cleanly when candidate pool is empty (0 candidates)", func(t *testing.T) {
		caseReader, candidateGen, candidateRanker, runRepo, _, uc := setup()

		run, _ := domain.NewMatchingRun("mrun_empty", "case_001", rankingVer, now)
		runRepo.FindByIDFn = func(ctx context.Context, id string) (*domain.MatchingRun, error) {
			return &run, nil
		}
		caseReader.GetCaseDetailsFn = func(ctx context.Context, caseID string) (*outbound.CaseDetails, error) {
			return &outbound.CaseDetails{CaseID: caseID, Category: cat}, nil
		}
		candidateGen.GenerateCandidatesFn = func(ctx context.Context, criteria domain.CandidateGenerationCriteria) (domain.CandidatePool, error) {
			return domain.NewCandidatePool([]domain.CandidateProfile{}, 100)
		}
		candidateRanker.RankFn = func(ctx context.Context, req outbound.RankingRequest) ([]domain.RankedCandidate, error) {
			return []domain.RankedCandidate{}, nil
		}
		var savedRun *domain.MatchingRun
		runRepo.SaveFn = func(ctx context.Context, r *domain.MatchingRun) error {
			savedRun = r
			return nil
		}

		err := uc.Execute(ctx, "mrun_empty")
		require.NoError(t, err)
		require.NotNil(t, savedRun)
		require.Equal(t, domain.RunStatusCompleted, savedRun.Status())
		require.Empty(t, savedRun.Candidates())
	})
}
