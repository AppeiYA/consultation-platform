package dto

import (
	"fmt"

	usecase_dto "github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
)

type UpdateAvailabilityRequest struct {
	AvailabilityID string `json:"availability_id" validate:"required"`
	DayOfWeek int `json:"day_of_week" validate:"required,min=0,max=6"`
	StartTime string `json:"start_time" validate:"required,datetime=15:04"`
	EndTime string `json:"end_time" validate:"required,datetime=15:04"`
}

func (r *UpdateAvailabilityRequest) Validate() error {
	if r.AvailabilityID == "" {
		return fmt.Errorf("availability_id is required")
	}
	if r.DayOfWeek < 0 || r.DayOfWeek > 6 {
		return fmt.Errorf("Invalid day of the week")
	}
	return nil
}

func (r *UpdateAvailabilityRequest) ToUsecaseDTO() *usecase_dto.UpdateAvailabilityRequest {
	return &usecase_dto.UpdateAvailabilityRequest{
		AvailabilityID: r.AvailabilityID,
		DayOfWeek: r.DayOfWeek,
		StartTime: r.StartTime,
		EndTime: r.EndTime,
	}
}