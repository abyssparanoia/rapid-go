package repository

import (
	"context"
	"database/sql"

	"github.com/abyssparanoia/rapid-go/api/src/lib/util"
	"github.com/abyssparanoia/rapid-go/api/src/model"
)

type user struct {
	sql *sql.DB
}

func (r *user) Get(ctx context.Context, userID int64) (*model.User, error) {
	user := &model.User{
		ID:         7777,
		Name:       "山田太郎",
		AvatarPath: "https://google.api.storage",
		Sex:        "man",
		Enabled:    true,
		CreatedAt:  util.TimeNow(),
		UpdatedAt:  util.TimeNow(),
	}

	return user, nil
}

// NewUser ... ユーザーレポジトリを取得する
func NewUser(sql *sql.DB) User {
	return &user{
		sql: sql,
	}
}
