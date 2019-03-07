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

// func (s *user) Create(ctx context.Context, Name string, AvatarPath string, Sex string) error {
// 	now := util.TimeNow()
// 	user := &model.User{
// 		Name:       Name,
// 		AvatarPath: AvatarPath,
// 		Sex:        Sex,
// 		Enabled:    true,
// 		CreatedAt:  now,
// 		UpdatedAt:  now,
// 	}
// 	err := s.uRepo.Insert(ctx, user)
// 	if err != nil {
// 		log.Errorf(ctx, "s.uRepo.Insert: %s", err.Error())
// 		return err
// 	}
// 	return nil
// }

// NewUser ... 新しいユーザーサービスを取得する
func NewUser(uRepo repository.User) User {
	return &user{
		uRepo: uRepo,
	}
}
