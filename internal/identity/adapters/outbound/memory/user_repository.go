package memory

import (
	"context"
	"errors"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	users map[string]*domain.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	r.users[user.ID()] = user
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	if _, ok := r.users[user.ID()]; !ok {
		return ErrUserNotFound
	}
	r.users[user.ID()] = user
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	for _, user := range r.users {
		if user.Email().String() == email.String() {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	if _, ok := r.users[id]; !ok {
		return ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}