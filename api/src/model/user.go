package model

import (
	"time"
)

// User ... ユーザーモデル
type User struct {
	ID         int64     ``
	Name       string    ``
	AvatarPath string    ``
	sex        string    ``
	Enabled    bool      ``
	CreatedAt  time.Time ``
	UpdatedAt  time.Time ``
}
