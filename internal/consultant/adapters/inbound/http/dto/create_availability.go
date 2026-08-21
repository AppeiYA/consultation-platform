package dto

import (
	"fmt"

	usecase_dto "github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
)

type CreateAvailabilityRequest struct {
	DayOfWeek int `json:"day_of_week" validate:"required,min=0,max=6"`
	StartTime string `json:"start_time" validate:"required,datetime=15:04"`
	EndTime string `json:"end_time" validate:"required,datetime=15:04"`
}

func (r *CreateAvailabilityRequest) Validate() error {
	if r.DayOfWeek < 0 || r.DayOfWeek > 6 {
		return fmt.Errorf("Invalid day of the week")
	}

	return nil
}

func (r *CreateAvailabilityRequest) ToUsecaseDTO() *usecase_dto.CreateAvailabilityRequest {
	return &usecase_dto.CreateAvailabilityRequest{
		DayOfWeek: r.DayOfWeek,
		StartTime: r.StartTime,
		EndTime: r.EndTime,
	}
}