package postgres

import (
	"time"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
)

type UserModel struct {
	ID             string    `db:"id"`
	FirstName      string    `db:"first_name"`
	LastName       string    `db:"last_name"`
	Email          string    `db:"email"`
	PasswordHash   string    `db:"password_hash"`
	Role           string    `db:"role"`
	IsVerified     bool      `db:"is_verified"`
	IsDeleted      bool      `db:"is_deleted"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func NewUserModel(user *domain.User) UserModel {
	return UserModel{
		ID:        user.ID(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Email:     user.Email().String(),
		PasswordHash:  user.PasswordHash().String(),
		Role:      string(user.Role()),
		IsVerified: user.IsVerified(),
		IsDeleted: user.IsDeleted(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}

func NewUserFromModel(model UserModel) (*domain.User, error) {
	return domain.ReconstituteUser(model.ID, model.FirstName, model.LastName, model.Email, model.PasswordHash, model.Role, model.IsDeleted, model.IsVerified, model.CreatedAt, model.UpdatedAt)
}
