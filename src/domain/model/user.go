package model

import "github.com/abyssparanoia/rapid-go/src/infrastructure/entity"

// User ... user model
type User struct {
	ID                  string  `json:"id"`
	DisplayName         string  `json:"display_name"`
	IconImagePath       string  `json:"icon_image_path"`
	BackgroundImagePath string  `json:"background_image_path"`
	Profile             *string `json:"profile"`
	Email               *string `json:"email"`
}

func newBaseUser(e *entity.User) *User {
	return &User{
		ID:                  e.ID,
		DisplayName:         e.DisplayName,
		IconImagePath:       e.IconImagePath,
		BackgroundImagePath: e.BackgroundImagePath,
		Profile:             e.Profile,
		Email:               e.Email,
	}
}

// NewUsers ... convert from entity to model
func NewUsers(users []*entity.User) []*User {
	dsts := []*User{}
	for _, user := range users {
		dsts = append(dsts, newBaseUser(user))
	}
	return dsts
}
