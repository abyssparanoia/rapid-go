package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/api/src/lib/log"
	"github.com/abyssparanoia/rapid-go/api/src/model"
	"github.com/jinzhu/gorm"
)

type user struct {
	sql *gorm.DB
}

func (r *user) Get(ctx context.Context, userID int64) (*model.User, error) {

	user := &model.User{}
	user.ID = userID
	errs := r.sql.First(user).GetErrors()
	if len(errs) != 0 {
		err := errs[0]
		log.Errorf(ctx, "r.sql.Find error: %s", err.Error())
		return nil, err
	}

	return user, nil
}

// NewUser ... ユーザーレポジトリを取得する
func NewUser(sql *gorm.DB) User {
	return &user{
		sql: sql,
	}
}
