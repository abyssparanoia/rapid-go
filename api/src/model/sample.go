package model

import (
	"time"
)

// Sample ... サンプルモデル
type Sample struct {
	ID        int64     `datastore:"-" boom:"id"`
	Category  string    ``
	Name      string    `datastore:",noindex"`
	Enabled   bool      ``
	CreatedAt time.Time ``
	UpdatedAt time.Time `datastore:"-"`
}
