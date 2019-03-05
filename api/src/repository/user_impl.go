package repository

import (
	"context"
	"database/sql"

	"github.com/abyssparanoia/rapid-go/api/src/lib/log"
	"github.com/abyssparanoia/rapid-go/api/src/lib/mysql"
	"github.com/abyssparanoia/rapid-go/api/src/model"

	sq "github.com/Masterminds/squirrel"
)

type user struct {
	sql *sql.DB
}

func (r *user) Get(ctx context.Context, userID int64) (*model.User, error) {

	q := sq.Select("id", "name", "avatar_path", "sex", "enabled", "created_at", "updated_at").
		From("m_users").
		Where(sq.Eq{"id": userID})

	mysql.DumpSelectQuery(ctx, q)

	rows, err := q.RunWith(r.sql).QueryContext(ctx)
	if err != nil {
		log.Errorf(ctx, "Get: %s", err.Error())
		return nil, err
	}

	var ret *model.User
	for rows.Next() {
		err := rows.Scan(
			&ret.ID,
			&ret.Name,
			&ret.AvatarPath,
			&ret.Sex,
			&ret.Enabled,
			&ret.CreatedAt,
			&ret.UpdatedAt)
		if err != nil {
			log.Errorf(ctx, "Get: %s", err.Error())
			rows.Close()
			return nil, err
		}
		break
	}

	return ret, nil
}

func (r *user) Insert(ctx context.Context, src *model.User) error {
	q := sq.Insert("m_users").
		Columns("name", "avatar_path", "sex", "enabled", "created_at", "updated_at").
		Values(src.Name, src.AvatarPath, src.Sex, src.Enabled, src.CreatedAt, src.UpdatedAt)

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
