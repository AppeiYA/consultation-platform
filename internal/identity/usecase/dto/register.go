package dto

type RegisterUserRequest struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type RegisterUserResponse struct {
	UserID string
}