package service

import (
	"context"

	"github.com/abyssparanoia/rapid-go/api/src/lib/log"
	"github.com/abyssparanoia/rapid-go/api/src/model"
	"github.com/abyssparanoia/rapid-go/api/src/repository"
)

type user struct {
	uRepo repository.User
}

func (s *user) Get(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.uRepo.Get(ctx, userID)
	if err != nil {
		log.Errorf(ctx, "s.uRepo.Get: %s", err.Error())
		return nil, err
	}
	return user, nil
}

// NewUser ... ユーザーサービスを取得する
func NewUser(uRepo repository.User) User {
	return &user{
		uRepo: uRepo,
	}
}
