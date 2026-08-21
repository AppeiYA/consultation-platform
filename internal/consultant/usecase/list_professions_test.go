package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/mocks"
)

func TestListProfessions_Execute(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	prof1 := domain.NewProfession("prof_9ee432d7-b672-40ae-b03f-c1f1fb696621", "SOFTWARE_ENGINEER", now)
	prof2 := domain.NewProfession("prof_12d965f5-e1f5-49aa-ac57-856772d236ce", "LAWYER", now)

	t.Run("should list professions successfully when repository returns professions", func(t *testing.T) {
		repo := &mocks.MockProfessionRepository{
			GetAllProfessionsFn: func(ctx context.Context) ([]*domain.Profession, error) {
				return []*domain.Profession{&prof1, &prof2}, nil
			},
		}
		sut := NewListProfessionsUsecase(repo)

		result, err := sut.Execute(context.Background())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 professions, got %d", len(result))
		}
		if result[0].ID() != "prof_9ee432d7-b672-40ae-b03f-c1f1fb696621" || result[0].Name() != "SOFTWARE_ENGINEER" {
			t.Errorf("unexpected first profession: %+v", result[0])
		}
		if result[1].ID() != "prof_12d965f5-e1f5-49aa-ac57-856772d236ce" || result[1].Name() != "LAWYER" {
			t.Errorf("unexpected second profession: %+v", result[1])
		}
	})

	t.Run("should return empty list when no professions exist in repository", func(t *testing.T) {
		repo := &mocks.MockProfessionRepository{
			GetAllProfessionsFn: func(ctx context.Context) ([]*domain.Profession, error) {
				return []*domain.Profession{}, nil
			},
		}
		sut := NewListProfessionsUsecase(repo)

		result, err := sut.Execute(context.Background())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 0 {
			t.Fatalf("expected 0 professions, got %d", len(result))
		}
	})

	t.Run("should fail when profession repository GetAllProfessions returns error", func(t *testing.T) {
		expectedErr := errors.New("database connection failed")
		repo := &mocks.MockProfessionRepository{
			GetAllProfessionsFn: func(ctx context.Context) ([]*domain.Profession, error) {
				return nil, expectedErr
			},
		}
		sut := NewListProfessionsUsecase(repo)

		result, err := sut.Execute(context.Background())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if result != nil {
			t.Errorf("expected nil result, got %v", result)
		}
	})
}
