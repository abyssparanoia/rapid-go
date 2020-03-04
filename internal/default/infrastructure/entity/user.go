package entity

import "github.com/abyssparanoia/rapid-go/internal/default/domain/model"

// User ... user entity
type User struct {
	ID                  string
	DisplayName         string
	IconImagePath       string
	BackgroundImagePath string
	Profile             *string
	Email               *string
	BaseEntity
}

// OutputModel ... output model from entity
func (e *User) OutputModel() *model.User {
	return &model.User{
		ID:                  e.ID,
		DisplayName:         e.DisplayName,
		IconImagePath:       e.IconImagePath,
		BackgroundImagePath: e.BackgroundImagePath,
		Profile:             e.Profile,
		Email:               e.Email,
	}
}

// NewUserFromModel ... new user from model
func NewUserFromModel(m *model.User) *User {
	e := &User{}
	e.ID = m.ID
	e.DisplayName = m.DisplayName
	e.IconImagePath = m.IconImagePath
	e.BackgroundImagePath = m.BackgroundImagePath
	e.Profile = m.Profile
	e.Email = m.Email
	return e
}

// OutputUsers ... output multi user models from entities
func OutputUsers(dsts []*User) []*model.User {

	users := []*model.User{}
	for _, dst := range dsts {
		users = append(users, dst.OutputModel())
	}

	return users
}

// UserTableName ... get user table name
func UserTableName() string {
	return "users"
}
