package dto

import usecase_dto "github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l LoginRequest) ToLoginParams() *usecase_dto.LoginRequest {
	return &usecase_dto.LoginRequest{
		Email:    l.Email,
		Password: l.Password,
	}
}