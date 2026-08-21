package usecase

import (
	"context"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

type CreateAvailabilityUsecase struct {
	availabilityRepository outbound.AvailabilityRepository
	idGenerator shared_outbound.IdentifierGenerator
	consultantRepository outbound.ConsultantRepository
	clock shared_outbound.Clock
}

func NewCreateAvailabilityUsecase(
	availabilityRepository outbound.AvailabilityRepository,
	idGenerator shared_outbound.IdentifierGenerator,
	consultantRepository outbound.ConsultantRepository,
	clock shared_outbound.Clock,
) *CreateAvailabilityUsecase {
	return &CreateAvailabilityUsecase{
		availabilityRepository: availabilityRepository,
		idGenerator: idGenerator,
		consultantRepository: consultantRepository,
		clock: clock,
	}
}

func (uc *CreateAvailabilityUsecase) Execute(ctx context.Context, userID string, req *dto.CreateAvailabilityRequest) error {
	// get consultant account by userID
	consultant, err := uc.consultantRepository.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if !consultant.IsAcceptingClients() {
		return domain.ErrConsultantNotAcceptingClients
	}

	// create TimeOfDay for start and end time
	startTime, err := domain.NewTimeOfDayFromString(req.StartTime)
	if err != nil {
		return err
	}

	endTime, err := domain.NewTimeOfDayFromString(req.EndTime)
	if err != nil {
		return err
	}

	// get consultant availability by day of week
	availabilities, err := uc.availabilityRepository.FindAvailabilitiesByConsultantIDAndDayOfWeek(ctx, consultant.ID(), time.Weekday(req.DayOfWeek))
	if err != nil {
		return err
	}

	// check for overlapping availability
	for _, availability := range availabilities {
		if availability.Overlaps(startTime, endTime) {
			return domain.ErrAvailabilityOverlap
		}
	}

	// generate ID for new availability
	availabilityID, err := uc.idGenerator.Generate(domain.ConsultantAvailabilityIDPrefix)
	if err != nil {
		return err
	}

	// create new consultant availability
	newAvailability, err := domain.NewConsultantAvailability(
		availabilityID,
		consultant.ID(),
		time.Weekday(req.DayOfWeek),
		startTime,
		endTime,
		uc.clock.Now(),
	)
	if err != nil {
		return err
	}

	return uc.availabilityRepository.SaveAvailability(ctx, newAvailability)
}