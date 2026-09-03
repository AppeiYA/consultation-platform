package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/mocks"
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/usecase/dto"
	shared_mocks "github.com/AppeiYA/consultation-platform/internal/shared/mocks"
	"github.com/stretchr/testify/require"
)

type testCreateCase struct {
	caseRepo        *mocks.MockCaseRepository
	clientVerifier  *mocks.MockClientVerifier
	matchingStarter *mocks.MockMatchingStarter
	idGenerator     *shared_mocks.MockIDGenerator
	clock           *shared_mocks.MockClock

	sut *CreateCaseUsecase
}

func setUpCreateCase(t *testing.T) *testCreateCase {
	t.Helper()

	caseRepo := &mocks.MockCaseRepository{
		SaveCaseFn: func(ctx context.Context, c *domain.ConsultationCase) error {
			return nil
		},
	}
	clientVerifier := &mocks.MockClientVerifier{
		VerifyClientFn: func(ctx context.Context, clientID string) error {
			return nil
		},
	}
	matchingStarter := &mocks.MockMatchingStarter{
		StartMatchingFn: func(ctx context.Context, caseID string) error {
			return nil
		},
	}
	idGenerator := &shared_mocks.MockIDGenerator{
		GenerateFn: func(prefix string) (string, error) {
			return "case_123456", nil
		},
	}
	clock := &shared_mocks.MockClock{
		NowFn: func() time.Time {
			return time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
		},
	}

	sut := NewCreateCaseUsecase(caseRepo, idGenerator, clientVerifier, matchingStarter, clock)

	return &testCreateCase{
		caseRepo:        caseRepo,
		clientVerifier:  clientVerifier,
		matchingStarter: matchingStarter,
		idGenerator:     idGenerator,
		clock:           clock,
		sut:             sut,
	}
}

func TestCreateCaseUsecase(t *testing.T) {
	ctx := context.Background()
	validReq := &dto.CreateCaseDTO{
		Title:       "Valid Case Title",
		Description: "Detailed description of the consultation case.",
		Category:    "FINANCE",
	}

	t.Run("should successfully create and save consultation case", func(t *testing.T) {
		tc := setUpCreateCase(t)
		var savedCase *domain.ConsultationCase
		tc.caseRepo.SaveCaseFn = func(ctx context.Context, c *domain.ConsultationCase) error {
			savedCase = c
			return nil
		}

		err := tc.sut.Execute(ctx, "user_123", validReq)
		require.NoError(t, err)
		require.NotNil(t, savedCase)
		require.Equal(t, "case_123456", savedCase.ID())
		require.Equal(t, "user_123", savedCase.ClientID())
		require.Equal(t, "Valid Case Title", savedCase.Title().String())
		require.Equal(t, "Detailed description of the consultation case.", savedCase.Description().String())
		require.Equal(t, "FINANCE", savedCase.Category().String())
		require.Equal(t, domain.CaseStatusDraft, savedCase.Status())
	})

	t.Run("should fail when client verification fails", func(t *testing.T) {
		tc := setUpCreateCase(t)
		clientErr := errors.New("client not found")
		tc.clientVerifier.VerifyClientFn = func(ctx context.Context, clientID string) error {
			return clientErr
		}

		err := tc.sut.Execute(ctx, "user_unknown", validReq)
		require.Error(t, err)
		require.Equal(t, clientErr, err)
	})

	t.Run("should fail when id generation fails", func(t *testing.T) {
		tc := setUpCreateCase(t)
		idErr := errors.New("id generator error")
		tc.idGenerator.GenerateFn = func(prefix string) (string, error) {
			return "", idErr
		}

		err := tc.sut.Execute(ctx, "user_123", validReq)
		require.Error(t, err)
		require.Equal(t, idErr, err)
	})

	t.Run("should fail when case title is invalid", func(t *testing.T) {
		tc := setUpCreateCase(t)
		badReq := &dto.CreateCaseDTO{
			Title:       "",
			Description: "Description",
			Category:    "FINANCE",
		}

		err := tc.sut.Execute(ctx, "user_123", badReq)
		require.Error(t, err)
		require.Equal(t, domain.ErrCaseTitleEmpty, err)
	})

	t.Run("should fail when case description is invalid", func(t *testing.T) {
		tc := setUpCreateCase(t)
		badReq := &dto.CreateCaseDTO{
			Title:       "Valid Title",
			Description: "",
			Category:    "FINANCE",
		}

		err := tc.sut.Execute(ctx, "user_123", badReq)
		require.Error(t, err)
		require.Equal(t, domain.ErrCaseDescriptionEmpty, err)
	})

	t.Run("should fail when case category is invalid", func(t *testing.T) {
		tc := setUpCreateCase(t)
		badReq := &dto.CreateCaseDTO{
			Title:       "Valid Title",
			Description: "Valid Description",
			Category:    "",
		}

		err := tc.sut.Execute(ctx, "user_123", badReq)
		require.Error(t, err)
		require.Equal(t, domain.ErrCaseCategoryEmpty, err)
	})

	t.Run("should fail when repository save fails", func(t *testing.T) {
		tc := setUpCreateCase(t)
		dbErr := errors.New("db save error")
		tc.caseRepo.SaveCaseFn = func(ctx context.Context, c *domain.ConsultationCase) error {
			return dbErr
		}

		err := tc.sut.Execute(ctx, "user_123", validReq)
		require.Error(t, err)
		require.Equal(t, dbErr, err)
	})

	t.Run("should auto-trigger matching when case is created", func(t *testing.T) {
		tc := setUpCreateCase(t)
		var triggeredCaseID string
		tc.matchingStarter.StartMatchingFn = func(ctx context.Context, caseID string) error {
			triggeredCaseID = caseID
			return nil
		}

		err := tc.sut.Execute(ctx, "user_123", validReq)
		require.NoError(t, err)
		require.Equal(t, "case_123456", triggeredCaseID)
	})
}
