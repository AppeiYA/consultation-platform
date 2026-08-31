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

func TestStartMatchingUsecase(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 8, 26, 12, 0, 0, 0, time.UTC)
	cat, _ := domain.NewMatchingCategory("Software Engineering")

	setup := func() (
		*mocks.MockCaseReader,
		*mocks.MockMatchingRunRepository,
		*mocks.MockMatchingJobEnqueuer,
		*shared_mocks.MockIDGenerator,
		*shared_mocks.MockClock,
		*StartMatchingUsecase,
	) {
		caseReader := &mocks.MockCaseReader{}
		runRepo := &mocks.MockMatchingRunRepository{}
		jobEnqueuer := &mocks.MockMatchingJobEnqueuer{}
		idGen := &shared_mocks.MockIDGenerator{}
		clock := &shared_mocks.MockClock{
			NowFn: func() time.Time { return now },
		}

		uc := NewStartMatchingUsecase(
			caseReader,
			runRepo,
			jobEnqueuer,
			idGen,
			clock,
		)
		return caseReader, runRepo, jobEnqueuer, idGen, clock, uc
	}

	t.Run("successful initiation of matching run and enqueuing job", func(t *testing.T) {
		caseReader, runRepo, jobEnqueuer, idGen, _, uc := setup()

		caseReader.GetCaseDetailsFn = func(ctx context.Context, caseID string) (*outbound.CaseDetails, error) {
			return &outbound.CaseDetails{
				CaseID:      caseID,
				ClientID:    "client_1",
				Category:    cat,
				Title:       "Backend Architecture Review",
				Description: "Need review for microservices",
			}, nil
		}

		idGen.GenerateFn = func(prefix string) (string, error) {
			return "mrun_abc123", nil
		}

		var savedRun *domain.MatchingRun
		runRepo.SaveFn = func(ctx context.Context, run *domain.MatchingRun) error {
			savedRun = run
			return nil
		}

		var enqueuedJob *outbound.MatchingJob
		jobEnqueuer.EnqueueFn = func(ctx context.Context, job outbound.MatchingJob) error {
			enqueuedJob = &job
			return nil
		}

		result, err := uc.Execute(ctx, "case_001")

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, "mrun_abc123", result.ID())
		require.Equal(t, "case_001", result.CaseID())
		require.Equal(t, domain.RunStatusPending, result.Status())
		require.Equal(t, savedRun, result)

		require.NotNil(t, enqueuedJob)
		require.Equal(t, "mrun_abc123", enqueuedJob.RunID)
		require.Equal(t, "case_001", enqueuedJob.CaseID)
	})

	t.Run("should reject empty case ID", func(t *testing.T) {
		_, _, _, _, _, uc := setup()
		_, err := uc.Execute(ctx, "")
		require.Equal(t, domain.ErrInvalidCaseID, err)
	})

	t.Run("should fail when case reader returns error (case not found)", func(t *testing.T) {
		caseReader, _, _, _, _, uc := setup()
		caseReader.GetCaseDetailsFn = func(ctx context.Context, caseID string) (*outbound.CaseDetails, error) {
			return nil, errors.New("case not found")
		}

		_, err := uc.Execute(ctx, "case_404")
		require.Error(t, err)
		require.Equal(t, "case not found", err.Error())
	})

	t.Run("should record failure and return error when job enqueuer fails", func(t *testing.T) {
		caseReader, runRepo, jobEnqueuer, idGen, _, uc := setup()

		caseReader.GetCaseDetailsFn = func(ctx context.Context, caseID string) (*outbound.CaseDetails, error) {
			return &outbound.CaseDetails{CaseID: caseID, Category: cat}, nil
		}
		idGen.GenerateFn = func(prefix string) (string, error) {
			return "mrun_fail_queue", nil
		}
		runRepo.SaveFn = func(ctx context.Context, run *domain.MatchingRun) error {
			return nil
		}
		jobEnqueuer.EnqueueFn = func(ctx context.Context, job outbound.MatchingJob) error {
			return errors.New("redis connection timeout")
		}

		_, err := uc.Execute(ctx, "case_001")
		require.Error(t, err)
		require.Equal(t, "redis connection timeout", err.Error())
	})
}
