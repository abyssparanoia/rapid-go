package model

import "github.com/abyssparanoia/rapid-go/src/infrastructure/entity"

// User ... ユーザーモデル
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Sex  string `json:"sex"`
}

// NewUserFromEntity ... entityからdomain modelへの変換をかねる
func NewUserFromEntity(e *entity.User) *User {
	return &User{
		ID:   e.ID,
		Name: e.Name,
	}
}
