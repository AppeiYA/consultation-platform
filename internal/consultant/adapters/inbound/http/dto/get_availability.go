package dto

import "github.com/AppeiYA/consultation-platform/internal/consultant/domain"

type GetAvailabilityResponse struct {
	DayOfWeek int `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func AvailabilityFromDomain(availability *domain.ConsultantAvailability) *GetAvailabilityResponse {
	return &GetAvailabilityResponse{
		DayOfWeek: int(availability.DayOfWeek()),
		StartTime: availability.StartTime().String(),
		EndTime:   availability.EndTime().String(),
	}
}