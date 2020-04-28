package model

import "time"

// User ... user model
type User struct {
	ID                  string    `json:"id"`
	DisplayName         string    `json:"display_name"`
	IconImagePath       string    `json:"icon_image_path"`
	BackgroundImagePath string    `json:"background_image_path"`
	Profile             *string   `json:"profile"`
	Email               *string   `json:"email"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// Exist ... check exist or not
func (user *User) Exist() bool {
	return user != nil
}

func NewUser(
	id,
	displayName,
	iconImagePath,
	backgroundImagePath,
	profile,
	email string,
) *User {

	return &User{
		ID:                  id,
		DisplayName:         displayName,
		IconImagePath:       iconImagePath,
		BackgroundImagePath: backgroundImagePath,
		Profile:             &profile,
		Email:               &email,
		CreatedAt:           time.Time{},
		UpdatedAt:           time.Time{},
	}
}
