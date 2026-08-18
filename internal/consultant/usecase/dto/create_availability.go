package dto

type CreateAvailabilityRequest struct {
	DayOfWeek    int  
	StartTime    string 
	EndTime      string 
}