package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/mocks"
	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
	shared_mocks "github.com/AppeiYA/consultation-platform/internal/shared/mocks"
	"github.com/stretchr/testify/require"
)

func TestExpertiseUsecases(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	consultant, _ := domain.ReconstitueConsultant(
		"con_123",
		"user_123",
		"prof_1",
		"John Expert",
		"Bio here",
		8,
		true,
		now,
		now,
	)

	setup := func() (
		*mocks.MockConsultantRepository,
		*mocks.MockExpertiseRepository,
		*shared_mocks.MockIDGenerator,
	) {
		consultantRepo := &mocks.MockConsultantRepository{
			FindByUserIDFn: func(ctx context.Context, userID string) (*domain.Consultant, error) {
				return consultant, nil
			},
		}
		expertiseRepo := &mocks.MockExpertiseRepository{}
		idGen := &shared_mocks.MockIDGenerator{
			GenerateFn: func(prefix string) (string, error) {
				return prefix + "_abc123", nil
			},
		}
		return consultantRepo, expertiseRepo, idGen
	}

	t.Run("AddExpertise success", func(t *testing.T) {
		consultantRepo, expertiseRepo, idGen := setup()
		uc := NewAddExpertiseUsecase(consultantRepo, expertiseRepo, idGen)

		var savedExp *domain.Expertise
		expertiseRepo.AddFn = func(ctx context.Context, expertise *domain.Expertise) error {
			savedExp = expertise
			return nil
		}

		res, err := uc.Execute(ctx, "user_123", dto.AddExpertiseDTO{Name: "Distributed Systems"})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "Distributed Systems", res.Name)
		require.Equal(t, "con_123", res.ConsultantID)
		require.Equal(t, "exp_abc123", res.ID)
		require.Equal(t, savedExp.Name(), res.Name)
	})

	t.Run("AddExpertise rejects empty name", func(t *testing.T) {
		consultantRepo, expertiseRepo, idGen := setup()
		uc := NewAddExpertiseUsecase(consultantRepo, expertiseRepo, idGen)

		_, err := uc.Execute(ctx, "user_123", dto.AddExpertiseDTO{Name: "   "})
		require.Equal(t, domain.ErrInvalidExpertiseName, err)
	})

	t.Run("RemoveExpertise success", func(t *testing.T) {
		consultantRepo, expertiseRepo, _ := setup()
		uc := NewRemoveExpertiseUsecase(consultantRepo, expertiseRepo)

		var deletedID string
		expertiseRepo.DeleteFn = func(ctx context.Context, consultantID string, expertiseID string) error {
			require.Equal(t, "con_123", consultantID)
			deletedID = expertiseID
			return nil
		}

		err := uc.Execute(ctx, "user_123", "exp_100")
		require.NoError(t, err)
		require.Equal(t, "exp_100", deletedID)
	})

	t.Run("ReplaceExpertises success", func(t *testing.T) {
		consultantRepo, expertiseRepo, idGen := setup()
		uc := NewReplaceExpertisesUsecase(consultantRepo, expertiseRepo, idGen)

		var replaced []*domain.Expertise
		expertiseRepo.ReplaceAllFn = func(ctx context.Context, consultantID string, expertises []*domain.Expertise) error {
			require.Equal(t, "con_123", consultantID)
			replaced = expertises
			return nil
		}

		res, err := uc.Execute(ctx, "user_123", dto.ReplaceExpertisesDTO{
			Expertises: []string{"Go", "PostgreSQL", "OAuth"},
		})
		require.NoError(t, err)
		require.Len(t, res, 3)
		require.Len(t, replaced, 3)
		require.Equal(t, "Go", res[0].Name)
		require.Equal(t, "PostgreSQL", res[1].Name)
		require.Equal(t, "OAuth", res[2].Name)
	})

	t.Run("ListMyExpertises success", func(t *testing.T) {
		consultantRepo, expertiseRepo, _ := setup()
		uc := NewListMyExpertisesUsecase(consultantRepo, expertiseRepo)

		exp1, _ := domain.NewExpertise("exp_1", "con_123", "Go")
		exp2, _ := domain.NewExpertise("exp_2", "con_123", "React")

		expertiseRepo.FindByConsultantIDFn = func(ctx context.Context, consultantID string) ([]*domain.Expertise, error) {
			require.Equal(t, "con_123", consultantID)
			return []*domain.Expertise{exp1, exp2}, nil
		}

		res, err := uc.Execute(ctx, "user_123")
		require.NoError(t, err)
		require.Len(t, res, 2)
		require.Equal(t, "Go", res[0].Name)
		require.Equal(t, "React", res[1].Name)
	})
}
