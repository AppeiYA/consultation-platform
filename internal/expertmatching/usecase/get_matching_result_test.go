package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/mocks"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/inbound"
	"github.com/stretchr/testify/require"
)

func TestGetMatchingResultUsecase(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 8, 26, 12, 0, 0, 0, time.UTC)
	v1, _ := domain.NewRankingVersion("v1")

	setupCompletedRun := func(numCandidates int) *domain.MatchingRun {
		run, _ := domain.NewMatchingRun("mrun_test", "case_test", v1, now)
		_ = run.StartGeneration()
		_ = run.StartRanking()

		candidates := make([]domain.RankedCandidate, numCandidates)
		for i := 0; i < numCandidates; i++ {
			r, _ := domain.NewRank(i + 1)
			s, _ := domain.NewMatchScore(float64(10-i) / 10.0)
			c, _ := domain.NewRankedCandidate("con_"+string(rune('A'+i)), r, s, nil)
			candidates[i] = c
		}
		_ = run.Complete(candidates, now)
		return &run
	}

	t.Run("should retrieve latest run by case ID and apply default Top 5", func(t *testing.T) {
		repo := &mocks.MockMatchingRunRepository{}
		run := setupCompletedRun(8)

		repo.FindLatestByCaseIDFn = func(ctx context.Context, caseID string) (*domain.MatchingRun, error) {
			require.Equal(t, "case_test", caseID)
			return run, nil
		}

		uc := NewGetMatchingResultUsecase(repo)
		resp, err := uc.Execute(ctx, inbound.GetMatchingResultRequest{CaseID: "case_test"})

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "mrun_test", resp.RunID)
		require.Equal(t, "case_test", resp.CaseID)
		require.Equal(t, "completed", resp.Status)
		require.Equal(t, "v1", resp.RankingVersion)
		require.Equal(t, 8, resp.TotalRanked)
		require.Len(t, resp.TopCandidates, 5) // default top 5
		require.Equal(t, "con_A", resp.TopCandidates[0].ConsultantID())
		require.Equal(t, "con_E", resp.TopCandidates[4].ConsultantID())
	})

	t.Run("should retrieve specific run by run ID with custom TopN", func(t *testing.T) {
		repo := &mocks.MockMatchingRunRepository{}
		run := setupCompletedRun(8)

		repo.FindByIDFn = func(ctx context.Context, id string) (*domain.MatchingRun, error) {
			require.Equal(t, "mrun_specific", id)
			return run, nil
		}

		uc := NewGetMatchingResultUsecase(repo)
		resp, err := uc.Execute(ctx, inbound.GetMatchingResultRequest{RunID: "mrun_specific", TopN: 3})

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.TopCandidates, 3)
		require.Equal(t, "con_C", resp.TopCandidates[2].ConsultantID())
	})

	t.Run("should reject when both CaseID and RunID are empty", func(t *testing.T) {
		repo := &mocks.MockMatchingRunRepository{}
		uc := NewGetMatchingResultUsecase(repo)

		_, err := uc.Execute(ctx, inbound.GetMatchingResultRequest{})
		require.Equal(t, domain.ErrInvalidCaseID, err)
	})

	t.Run("should return not found error when run does not exist", func(t *testing.T) {
		repo := &mocks.MockMatchingRunRepository{}
		repo.FindLatestByCaseIDFn = func(ctx context.Context, caseID string) (*domain.MatchingRun, error) {
			return nil, nil
		}

		uc := NewGetMatchingResultUsecase(repo)
		_, err := uc.Execute(ctx, inbound.GetMatchingResultRequest{CaseID: "case_missing"})
		require.Error(t, err)
		require.Equal(t, "matching run not found", err.Error())
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		repo := &mocks.MockMatchingRunRepository{}
		repo.FindLatestByCaseIDFn = func(ctx context.Context, caseID string) (*domain.MatchingRun, error) {
			return nil, errors.New("db connection failure")
		}

		uc := NewGetMatchingResultUsecase(repo)
		_, err := uc.Execute(ctx, inbound.GetMatchingResultRequest{CaseID: "case_001"})
		require.Error(t, err)
		require.Equal(t, "db connection failure", err.Error())
	})
}
