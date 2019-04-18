package entity

// User ... ユーザーエンティティ
type User struct {
	ID        int64
	Name      string
	Sex       string
	Enabled   bool
	CreatedAt int64
	UpdatedAt int64
}
