package dto

type RegisterUserRequest struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Role      string
}

type RegisterUserResponse struct {
	UserID string
}