package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
)

type GetAvailabilityUsecase struct {
	availabilityRepository outbound.AvailabilityRepository
}

func NewGetAvailabilityUsecase(
	availabilityRepository outbound.AvailabilityRepository,
) *GetAvailabilityUsecase {
	return &GetAvailabilityUsecase{
		availabilityRepository: availabilityRepository,
	}
}

func (uc *GetAvailabilityUsecase) Execute(ctx context.Context, consultantID string) ([]*domain.ConsultantAvailability, error) {
	// get consultant availability by consultantID
	availabilities, err := uc.availabilityRepository.FindAvailabilitiesByConsultantID(ctx, consultantID)
	if err != nil {
		return nil, err
	}
	
	return availabilities, nil
}