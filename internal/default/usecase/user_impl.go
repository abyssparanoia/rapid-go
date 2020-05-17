package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/default/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/default/domain/repository"
)

type user struct {
	userRepo repository.User
}

func (s *user) Get(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// NewUser ... get User usecase
func NewUser(userRepo repository.User) User {
	return &user{userRepo}
}
