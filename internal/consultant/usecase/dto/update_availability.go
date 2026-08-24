package dto

type UpdateAvailabilityRequest struct {
	AvailabilityID string 
	DayOfWeek int 
	StartTime string
	EndTime   string
}