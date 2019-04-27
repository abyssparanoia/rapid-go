package model

import "github.com/abyssparanoia/rapid-go/src/infrastructure/entity"

// User ... user model
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Sex  string `json:"sex"`
}

// NewUserFromEntity ... convert from entity to model
func NewUserFromEntity(e *entity.User) *User {
	return &User{
		ID:   e.ID,
		Name: e.Name,
		Sex:  e.Sex,
	}
}
