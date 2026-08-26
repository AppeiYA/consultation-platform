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

func TestDeleteCaseUsecase(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	title, _ := domain.NewCaseTitle("Test Case")
	desc, _ := domain.NewCaseDescription("Test Description")
	cat, _ := domain.NewCaseCategory("LEGAL")
	testCase, _ := domain.NewConsultationCase("case_123", "user_123", *title, *desc, *cat, now)

	t.Run("should delete consultation case successfully", func(t *testing.T) {
		deletedID := ""
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return testCase, nil
			},
			DeleteCaseFn: func(ctx context.Context, id string) error {
				deletedID = id
				return nil
			},
		}
		sut := NewDeleteCaseUsecase(repo)

		err := sut.Execute(ctx, "user_123", "case_123")
		require.NoError(t, err)
		require.Equal(t, "case_123", deletedID)
	})

	t.Run("should return ErrCaseNotFound when case does not exist", func(t *testing.T) {
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return nil, nil
			},
		}
		sut := NewDeleteCaseUsecase(repo)

		err := sut.Execute(ctx, "user_123", "case_nonexistent")
		require.Error(t, err)
		require.Equal(t, domain.ErrCaseNotFound, err)
	})

	t.Run("should return ErrUnauthorizedAccess when case belongs to another user", func(t *testing.T) {
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return testCase, nil
			},
		}
		sut := NewDeleteCaseUsecase(repo)

		err := sut.Execute(ctx, "user_different", "case_123")
		require.Error(t, err)
		require.Equal(t, domain.ErrUnauthorizedAccess, err)
	})

	t.Run("should return error when repository delete fails", func(t *testing.T) {
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return testCase, nil
			},
			DeleteCaseFn: func(ctx context.Context, id string) error {
				return errors.New("delete failed")
			},
		}
		sut := NewDeleteCaseUsecase(repo)

		err := sut.Execute(ctx, "user_123", "case_123")
		require.Error(t, err)
		require.Equal(t, "delete failed", err.Error())
	})
}
