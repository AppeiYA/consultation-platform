package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/mocks"
)

func TestGetProfession_Execute(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	profession := domain.NewProfession("prof_001", "SOFTWARE_ENGINEER", now)

	t.Run("should get profession successfully", func(t *testing.T) {
		repo := &mocks.MockProfessionRepository{
			GetProfessionByIDFn: func(ctx context.Context, professionID string) (*domain.Profession, error) {
				return &profession, nil
			},
		}
		sut := NewGetProfessionUsecase(repo)

		res, err := sut.Execute(context.Background(), "prof_001")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res == nil {
			t.Fatal("expected profession, got nil")
		}
		if res.ID() != "prof_001" {
			t.Errorf("expected ID prof_001, got %s", res.ID())
		}
		if res.Name() != "SOFTWARE_ENGINEER" {
			t.Errorf("expected Name SOFTWARE_ENGINEER, got %s", res.Name())
		}
	})

	t.Run("should return error when repo returns error", func(t *testing.T) {
		expectedErr := errors.New("db error")
		repo := &mocks.MockProfessionRepository{
			GetProfessionByIDFn: func(ctx context.Context, professionID string) (*domain.Profession, error) {
				return nil, expectedErr
			},
		}
		sut := NewGetProfessionUsecase(repo)

		res, err := sut.Execute(context.Background(), "prof_001")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("expected nil response, got %v", res)
		}
	})
}
