package model

import (
	"database/sql"
	"time"
)

// User ... ユーザーモデル
type User struct {
	ID         int64          `db:"id" json:"id"`
	Name       string         `db:"name" json:"name"`
	AvatarPath sql.NullString `db:"avatar_path" json:"avatar_path"`
	Sex        string         `db:"sex" json:"sex"`
	Enabled    bool           `db:"enabled" json:"enabled"`
	CreatedAt  time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at" json:"updated_at"`
}
