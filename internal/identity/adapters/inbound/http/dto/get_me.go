package dto

import usecase_dto "github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"

type GetMeResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

func NewGetMeResponse(user *usecase_dto.GetCurrentUserResponse) *GetMeResponse {
	return &GetMeResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
	}
}