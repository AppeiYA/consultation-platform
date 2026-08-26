package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/mocks"
	"github.com/stretchr/testify/require"
)

func TestGetCaseUsecase(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	title, _ := domain.NewCaseTitle("Test Case")
	desc, _ := domain.NewCaseDescription("Test Description")
	cat, _ := domain.NewCaseCategory("LEGAL")
	testCase, _ := domain.NewConsultationCase("case_123", "user_123", *title, *desc, *cat, now)

	t.Run("should retrieve consultation case successfully", func(t *testing.T) {
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return testCase, nil
			},
		}
		sut := NewGetCaseUsecase(repo)

		result, err := sut.Execute(ctx, "user_123", "case_123")
		require.NoError(t, err)
		require.Equal(t, testCase, result)
	})

	t.Run("should return ErrCaseNotFound when case is not found", func(t *testing.T) {
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return nil, nil
			},
		}
		sut := NewGetCaseUsecase(repo)

		result, err := sut.Execute(ctx, "user_123", "case_nonexistent")
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, domain.ErrCaseNotFound, err)
	})

	t.Run("should return ErrUnauthorizedAccess when case belongs to another user", func(t *testing.T) {
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return testCase, nil
			},
		}
		sut := NewGetCaseUsecase(repo)

		result, err := sut.Execute(ctx, "user_different", "case_123")
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, domain.ErrUnauthorizedAccess, err)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return nil, errors.New("db error")
			},
		}
		sut := NewGetCaseUsecase(repo)

		result, err := sut.Execute(ctx, "user_123", "case_123")
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, "db error", err.Error())
	})
}
