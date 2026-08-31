package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

type RegisterConsultant struct {
	consultantRepo  outbound.ConsultantRepository
	professionRepo  outbound.ProfessionRepository
	expertiseRepo   outbound.ExpertiseRepository
	identityAdapter outbound.RoleAssigner
	idGenerator     shared_outbound.IdentifierGenerator
	clock           shared_outbound.Clock
}

func NewRegisterConsultantUsecase(
	consultantRepo outbound.ConsultantRepository,
	professionRepo outbound.ProfessionRepository,
	expertiseRepo outbound.ExpertiseRepository,
	identityAdapter outbound.RoleAssigner,
	idGenerator shared_outbound.IdentifierGenerator,
	clock shared_outbound.Clock,
) *RegisterConsultant {
	return &RegisterConsultant{
		consultantRepo:  consultantRepo,
		professionRepo:  professionRepo,
		expertiseRepo:   expertiseRepo,
		identityAdapter: identityAdapter,
		idGenerator:     idGenerator,
		clock:           clock,
	}
}

func (uc *RegisterConsultant) Execute(ctx context.Context, userID string, req *dto.RegisterConsultantDTO) error {
	// check for existing consultant with the same user ID
	exists, err := uc.consultantRepo.ExistsByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrConsultantAlreadyExists
	}

	if req.ProfessionID == "" {
		return domain.ErrInvalidProfession
	}

	// check if profession exists
	profession, err := uc.professionRepo.GetProfessionByID(ctx, req.ProfessionID)
	if err != nil {
		return err
	}
	if profession == nil {
		return domain.ErrInvalidProfession
	}

	// generate a new ID for consultant
	newID, err := uc.idGenerator.Generate(domain.ConsultantIDPrefix)
	if err != nil {
		return err
	}

	// check correct display name
	displayName, err := domain.NewDisplayName(req.DisplayName)
	if err != nil {
		return err
	}

	// check correct bio
	bio, err := domain.NewBio(req.Bio)
	if err != nil {
		return err
	}

	yearsExperience, err := domain.NewYearsExperience(req.YearsExperience)
	if err != nil {
		return err
	}

	// create a new consultant
	consultant := domain.NewConsultant(
		newID,
		userID,
		profession.ID(),
		displayName,
		bio,
		yearsExperience,
		uc.clock.Now(),
	)

	err = uc.consultantRepo.Save(ctx, consultant)
	if err != nil {
		return err
	}

	if len(req.Expertises) > 0 {
		expertiseEntities := make([]*domain.Expertise, 0, len(req.Expertises))
		for _, name := range req.Expertises {
			expID, err := uc.idGenerator.Generate(domain.ExpertiseIDPrefix)
			if err != nil {
				return err
			}
			exp, err := domain.NewExpertise(expID, newID, name)
			if err != nil {
				return err
			}
			expertiseEntities = append(expertiseEntities, exp)
		}
		if err := uc.expertiseRepo.SaveMany(ctx, expertiseEntities); err != nil {
			return err
		}
	}

	return uc.identityAdapter.AssignConsultantRole(ctx, userID)
}
