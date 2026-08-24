package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

type ActivateAvailabilityUsecase struct {
	availabilityRepository outbound.AvailabilityRepository
	consultantRepository  outbound.ConsultantRepository
	clock shared_outbound.Clock
}

func NewActivateAvailabilityUsecase(
	availabilityRepository outbound.AvailabilityRepository,
	consultantRepository outbound.ConsultantRepository,
	clock shared_outbound.Clock,
) *ActivateAvailabilityUsecase {
	return &ActivateAvailabilityUsecase{
		availabilityRepository: availabilityRepository,
		consultantRepository: consultantRepository,
		clock: clock,
	}
}

func (uc *ActivateAvailabilityUsecase) Execute(ctx context.Context, userID string, availabilityID string) error {
	// get consultant account by userID
	consultant, err := uc.consultantRepository.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// get availability by ID
	availability, err := uc.availabilityRepository.FindAvailabilityByID(ctx, availabilityID)
	if err != nil {
		return err
	}
	if availability == nil {
		return domain.ErrAvailabilityNotFound
	}

	// confirm that the availability belongs to the consultant
	if availability.ConsultantID() != consultant.ID() {
		return domain.ErrAvailabilityNotFound
	}

	if availability.IsActive() {
		return domain.ErrAvailabilityAlreadyActivated
	}

	// activate the availability
	err = uc.availabilityRepository.ActivateAvailability(ctx, availabilityID)
	if err != nil {
		return err
	}

	return nil
}