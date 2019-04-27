package entity

// User ... user entity
type User struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
	Sex  string `db:"sex"`
	BaseEntity
}
