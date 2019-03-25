package model

// User ... ユーザーモデル
type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Sex       string `json:"sex"`
	Enabled   bool   `json:"enabled"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
