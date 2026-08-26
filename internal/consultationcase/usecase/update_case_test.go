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

func TestUpdateCaseUsecase(t *testing.T) {
	ctx := context.Background()
	initialTime := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	updateTime := time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)

	newTestCase := func() *domain.ConsultationCase {
		title, _ := domain.NewCaseTitle("Original Title")
		desc, _ := domain.NewCaseDescription("Original Description")
		cat, _ := domain.NewCaseCategory("ORIGINAL_CAT")
		c, _ := domain.NewConsultationCase("case_123", "user_123", *title, *desc, *cat, initialTime)
		return c
	}

	clock := &shared_mocks.MockClock{
		NowFn: func() time.Time {
			return updateTime
		},
	}

	t.Run("should update title, description, and category successfully", func(t *testing.T) {
		c := newTestCase()
		var updatedCase *domain.ConsultationCase
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return c, nil
			},
			UpdateCaseFn: func(ctx context.Context, c *domain.ConsultationCase) error {
				updatedCase = c
				return nil
			},
		}
		sut := NewUpdateCaseUsecase(repo, clock)

		newTitle := "Updated Title"
		newDesc := "Updated Description"
		newCat := "NEW_CAT"
		req := &dto.UpdateCaseDTO{
			Title:       &newTitle,
			Description: &newDesc,
			Category:    &newCat,
		}

		err := sut.Execute(ctx, "user_123", "case_123", req)
		require.NoError(t, err)
		require.NotNil(t, updatedCase)
		require.Equal(t, "Updated Title", updatedCase.Title().String())
		require.Equal(t, "Updated Description", updatedCase.Description().String())
		require.Equal(t, "NEW_CAT", updatedCase.Category().String())
		require.Equal(t, updateTime, updatedCase.UpdatedAt())
	})

	t.Run("should update only provided fields (partial update)", func(t *testing.T) {
		c := newTestCase()
		var updatedCase *domain.ConsultationCase
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return c, nil
			},
			UpdateCaseFn: func(ctx context.Context, c *domain.ConsultationCase) error {
				updatedCase = c
				return nil
			},
		}
		sut := NewUpdateCaseUsecase(repo, clock)

		newTitle := "Only Title Updated"
		req := &dto.UpdateCaseDTO{
			Title: &newTitle,
		}

		err := sut.Execute(ctx, "user_123", "case_123", req)
		require.NoError(t, err)
		require.NotNil(t, updatedCase)
		require.Equal(t, "Only Title Updated", updatedCase.Title().String())
		require.Equal(t, "Original Description", updatedCase.Description().String())
		require.Equal(t, "ORIGINAL_CAT", updatedCase.Category().String())
	})

	t.Run("should fail when case not found", func(t *testing.T) {
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return nil, nil
			},
		}
		sut := NewUpdateCaseUsecase(repo, clock)

		newTitle := "New Title"
		err := sut.Execute(ctx, "user_123", "case_unknown", &dto.UpdateCaseDTO{Title: &newTitle})
		require.Error(t, err)
		require.Equal(t, domain.ErrCaseNotFound, err)
	})

	t.Run("should fail when case belongs to another user", func(t *testing.T) {
		c := newTestCase()
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return c, nil
			},
		}
		sut := NewUpdateCaseUsecase(repo, clock)

		newTitle := "New Title"
		err := sut.Execute(ctx, "user_other", "case_123", &dto.UpdateCaseDTO{Title: &newTitle})
		require.Error(t, err)
		require.Equal(t, domain.ErrUnauthorizedAccess, err)
	})

	t.Run("should fail when new title is empty string", func(t *testing.T) {
		c := newTestCase()
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return c, nil
			},
		}
		sut := NewUpdateCaseUsecase(repo, clock)

		emptyTitle := ""
		err := sut.Execute(ctx, "user_123", "case_123", &dto.UpdateCaseDTO{Title: &emptyTitle})
		require.Error(t, err)
		require.Equal(t, domain.ErrCaseTitleEmpty, err)
	})

	t.Run("should fail when new description is empty string", func(t *testing.T) {
		c := newTestCase()
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return c, nil
			},
		}
		sut := NewUpdateCaseUsecase(repo, clock)

		emptyDesc := ""
		err := sut.Execute(ctx, "user_123", "case_123", &dto.UpdateCaseDTO{Description: &emptyDesc})
		require.Error(t, err)
		require.Equal(t, domain.ErrCaseDescriptionEmpty, err)
	})

	t.Run("should fail when new category is empty string", func(t *testing.T) {
		c := newTestCase()
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return c, nil
			},
		}
		sut := NewUpdateCaseUsecase(repo, clock)

		emptyCat := ""
		err := sut.Execute(ctx, "user_123", "case_123", &dto.UpdateCaseDTO{Category: &emptyCat})
		require.Error(t, err)
		require.Equal(t, domain.ErrCaseCategoryEmpty, err)
	})

	t.Run("should fail when repository update fails", func(t *testing.T) {
		c := newTestCase()
		repo := &mocks.MockCaseRepository{
			FindCaseByIDFn: func(ctx context.Context, id string) (*domain.ConsultationCase, error) {
				return c, nil
			},
			UpdateCaseFn: func(ctx context.Context, c *domain.ConsultationCase) error {
				return errors.New("db error")
			},
		}
		sut := NewUpdateCaseUsecase(repo, clock)

		newTitle := "New Title"
		err := sut.Execute(ctx, "user_123", "case_123", &dto.UpdateCaseDTO{Title: &newTitle})
		require.Error(t, err)
		require.Equal(t, "db error", err.Error())
	})
}
