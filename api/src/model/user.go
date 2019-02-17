package model

import (
	"time"
)

// User ... ユーザーモデル
type User struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	AvatarPath string    `json:"avatar_path"`
	Sex        string    `json:"sex"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
