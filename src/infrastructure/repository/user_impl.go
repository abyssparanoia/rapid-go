package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/src/domain/repository"
	"github.com/abyssparanoia/rapid-go/src/infrastructure/entity"
	"github.com/abyssparanoia/rapid-go/src/lib/log"
	"github.com/jmoiron/sqlx"
)

type user struct {
	sql *sqlx.DB
}

func (r *user) Get(ctx context.Context, userID int64) (*entity.User, error) {

	user := &entity.User{}
	user.ID = userID
	rows, err := r.sql.Queryx("SELECT * FROM users WHERE id=1")
	if err != nil {
		log.Errorf(ctx, "r.sql.Queryx: %s", err.Error())
		return nil, err
	}

	for rows.Next() {
		err := rows.StructScan(&user)
		if err != nil {
			log.Errorf(ctx, "r.sql.Queryx: %s", err.Error())
			return nil, err
		}
		break
	}

	return user, nil
}

// NewUser ... ユーザーレポジトリを取得する
func NewUser(sql *sqlx.DB) repository.User {
	return &user{
		sql: sql,
	}
}
