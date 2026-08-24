package usecase

import (
	"context"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

type UpdateAvailabilityUsecase struct {
	consultantRepository outbound.ConsultantRepository
	availabilityRepository outbound.AvailabilityRepository
	clock shared_outbound.Clock
}

func NewUpdateAvailabilityUsecase(
	availabilityRepository outbound.AvailabilityRepository, 
	consultantRepository outbound.ConsultantRepository, 
	clock shared_outbound.Clock,
	) *UpdateAvailabilityUsecase {
	return &UpdateAvailabilityUsecase{
		availabilityRepository: availabilityRepository,
		consultantRepository: consultantRepository,
		clock: clock,
	}
}	

func (u *UpdateAvailabilityUsecase) Execute(ctx context.Context, userID string, req *dto.UpdateAvailabilityRequest) (*domain.ConsultantAvailability, error) {
	// get consultant by userID
	consultant, err := u.consultantRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// get availability by ID
	availability, err := u.availabilityRepository.FindAvailabilityByID(ctx, req.AvailabilityID)
	if err != nil {
		return nil, err
	}
	if availability == nil {
		return nil, domain.ErrAvailabilityNotFound
	}

	// confirm that the availability belongs to the consultant
	if availability.ConsultantID() != consultant.ID() {
		return nil, domain.ErrAvailabilityNotFound
	}

	var startTime, endTime domain.TimeOfDay
	if req.StartTime != "" {
		startTime, err = domain.NewTimeOfDayFromString(req.StartTime)
		if err != nil {
			return nil, err
		}
	}else {
		startTime = availability.StartTime()
	}

	if req.EndTime != "" {
		endTime, err = domain.NewTimeOfDayFromString(req.EndTime)
		if err != nil {
			return nil, err
		}
	}else {
		endTime = availability.EndTime()
	}

	// check for overlapping availability with other slots on the same day
	availabilities, err := u.availabilityRepository.FindAvailabilitiesByConsultantIDAndDayOfWeek(ctx, consultant.ID(), time.Weekday(req.DayOfWeek))
	if err != nil {
		return nil, err
	}

	for _, existing := range availabilities {
		if existing.ID() != availability.ID() && existing.Overlaps(startTime, endTime) {
			return nil, domain.ErrAvailabilityOverlap
		}
	}

	err = availability.Update(
		time.Weekday(req.DayOfWeek), 
		startTime, 
		endTime, 
		u.clock.Now(),
	)
	if err != nil {
		return nil, err
	}

	err = u.availabilityRepository.UpdateAvailability(ctx, availability)
	if err != nil {
		return nil, err
	}

	return availability, nil 
}