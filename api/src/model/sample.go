package model

import (
	"time"
)

// Sample ... サンプルモデル
type Sample struct {
	ID        int64           `datastore:"-" boom:"id"`
	Category  string          ``
	Name      string          `datastore:",noindex"`
	Enabled   bool            ``
	Details   []*SampleDetail `datastore:",flatten,noindex"`
	CreatedAt time.Time       ``
	UpdatedAt time.Time       `datastore:"-"`
}

type SampleDetail struct {
	Name   string `datastore:""`
	Detail string `datastore:""`
}
