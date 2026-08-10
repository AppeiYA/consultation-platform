package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
)

type RegisterConsultant struct {
	consultantRepo outbound.ConsultantRepository
	idGenerator shared_outbound.IdentifierGenerator
	clock shared_outbound.Clock
}

func NewRegisterConsultantUsecase(consultantRepo outbound.ConsultantRepository, idGenerator shared_outbound.IdentifierGenerator, clock shared_outbound.Clock) *RegisterConsultant {
	return &RegisterConsultant{
		consultantRepo: consultantRepo,
		idGenerator: idGenerator,
		clock: clock,
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

	// genarate a new ID for consultant

	newID, err := uc.idGenerator.Generate(domain.ConsultantIDPrefix)
	if err != nil {
		return err
	}

	// check profession is valid
	profession, err := domain.NewProfession(req.Profession)
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
		profession,
		displayName,
		bio,
		yearsExperience,
		uc.clock.Now(),
	)

	return uc.consultantRepo.Save(ctx, consultant)
}