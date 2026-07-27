package dto

import usecase_dto "github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=50"`
	LastName  string `json:"last_name" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=50"`
	Role      string `json:"role" validate:"required,oneof=CLIENT EXPERT"`
}

type RegisterUserResponse struct {
	UserID string `json:"user_id"`
}

func NewRegisterUserResponse(userID string) RegisterUserResponse {
	return RegisterUserResponse{
		UserID: userID,
	}
}

func (r RegisterRequest) ToRegisterParams() *usecase_dto.RegisterUserRequest {
	return &usecase_dto.RegisterUserRequest{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Password:  r.Password,
		Role:      r.Role,
	}
}