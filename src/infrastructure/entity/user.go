package entity

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

// UserTableName ... get user table name
func UserTableName() string {
	return "users"
}
