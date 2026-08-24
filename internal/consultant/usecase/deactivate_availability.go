package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)
type DeactivateAvailabilityUsecase struct {
	availabilityRepository outbound.AvailabilityRepository
	consultantRepository  outbound.ConsultantRepository
	clock shared_outbound.Clock
}

func NewDeactivateAvailabilityUsecase(
	availabilityRepository outbound.AvailabilityRepository,
	consultantRepository outbound.ConsultantRepository,
	clock shared_outbound.Clock,
) *DeactivateAvailabilityUsecase {
	return &DeactivateAvailabilityUsecase{
		availabilityRepository: availabilityRepository,
		consultantRepository: consultantRepository,
		clock: clock,
	}
}

func (uc *DeactivateAvailabilityUsecase) Execute(ctx context.Context, userID string, availabilityID string) error {
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

	if !availability.IsActive() {
		return domain.ErrAvailabilityAlreadyDeactivated
	}

	// deactivate the availability
	err = uc.availabilityRepository.DeactivateAvailability(ctx, availabilityID)
	if err != nil {
		return err
	}

	return nil
}