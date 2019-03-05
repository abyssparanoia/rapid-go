package repository

import (
	"context"
	"database/sql"

	"github.com/abyssparanoia/rapid-go/api/src/lib/log"
	"github.com/abyssparanoia/rapid-go/api/src/lib/mysql"
	"github.com/abyssparanoia/rapid-go/api/src/lib/util"
	"github.com/abyssparanoia/rapid-go/api/src/model"

	sq "github.com/Masterminds/squirrel"
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

func (r *user) Insert(ctx context.Context, src *model.User) error {
	q := sq.Insert("m_users").
		Columns("id", "name", "avatar_path", "sex", "enabled", "created_at", "updated_at").
		Values(src.ID, src.Name, src.AvatarPath, src.Sex, src.Enabled, src.CreatedAt, src.UpdatedAt)

	mysql.DumpInsertQuery(ctx, q)

	_, err := q.RunWith(r.sql).ExecContext(ctx)
	if err != nil {
		log.Errorf(ctx, "MySQLInsert: %s", err.Error())
		return err
	}

	return nil
}

// NewUser ... ユーザーレポジトリを取得する
func NewUser(sql *sql.DB) User {
	return &user{
		sql: sql,
	}
}
