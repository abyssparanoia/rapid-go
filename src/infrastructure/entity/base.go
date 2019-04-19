package entity

// BaseEntity ... 共通の構造体実装
type BaseEntity struct {
	Enabled   bool  `db:"enabled"`
	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}
