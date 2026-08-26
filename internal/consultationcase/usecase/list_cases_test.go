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

func TestListCasesUsecase(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	title1, _ := domain.NewCaseTitle("Case 1")
	desc1, _ := domain.NewCaseDescription("Desc 1")
	cat1, _ := domain.NewCaseCategory("LEGAL")
	case1, _ := domain.NewConsultationCase("case_1", "user_123", *title1, *desc1, *cat1, now)

	title2, _ := domain.NewCaseTitle("Case 2")
	desc2, _ := domain.NewCaseDescription("Desc 2")
	cat2, _ := domain.NewCaseCategory("TECH")
	case2, _ := domain.NewConsultationCase("case_2", "user_123", *title2, *desc2, *cat2, now)

	t.Run("should list all consultation cases for client", func(t *testing.T) {
		repo := &mocks.MockCaseRepository{
			FindCasesByClientIDFn: func(ctx context.Context, clientID string) ([]*domain.ConsultationCase, error) {
				return []*domain.ConsultationCase{case1, case2}, nil
			},
		}
		sut := NewListCasesUsecase(repo)

		results, err := sut.Execute(ctx, "user_123")
		require.NoError(t, err)
		require.Len(t, results, 2)
		require.Equal(t, case1, results[0])
		require.Equal(t, case2, results[1])
	})

	t.Run("should return empty list when client has no cases", func(t *testing.T) {
		repo := &mocks.MockCaseRepository{
			FindCasesByClientIDFn: func(ctx context.Context, clientID string) ([]*domain.ConsultationCase, error) {
				return []*domain.ConsultationCase{}, nil
			},
		}
		sut := NewListCasesUsecase(repo)

		results, err := sut.Execute(ctx, "user_123")
		require.NoError(t, err)
		require.Empty(t, results)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		repo := &mocks.MockCaseRepository{
			FindCasesByClientIDFn: func(ctx context.Context, clientID string) ([]*domain.ConsultationCase, error) {
				return nil, errors.New("database failure")
			},
		}
		sut := NewListCasesUsecase(repo)

		results, err := sut.Execute(ctx, "user_123")
		require.Error(t, err)
		require.Nil(t, results)
		require.Equal(t, "database failure", err.Error())
	})
}
