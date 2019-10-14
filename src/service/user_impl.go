package service

import (
	"context"

	"github.com/abyssparanoia/rapid-go/src/domain/model"
	"github.com/abyssparanoia/rapid-go/src/domain/repository"
	"github.com/abyssparanoia/rapid-go/src/lib/log"
)

type user struct {
	userRepo repository.User
}

func (s *user) Get(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		log.Errorm(ctx, "s.userRepo.Get", err)
		return nil, err
	}
	return user, nil
}

// NewUser ... get User service
func NewUser(userRepo repository.User) User {
	return &user{userRepo}
}
