package service

import (
	"context"

	"github.com/abyssparanoia/rapid-go/api/src/model"
)

// User ... ユーザーサービスのインターフェイス
type User interface {
	Get(ctx context.Context, userID int64) (*model.User, error)
	Create(ctx context.Context, Name string, AvatarPath string, Sex string) error
}
