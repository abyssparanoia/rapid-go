package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/volatiletech/null/v8"
)

type User interface {
	Get(
		ctx context.Context,
		query GetUserQuery,
		orFail bool,
		preload bool,
	) (*model.User, error)
	Create(
		ctx context.Context,
		user *model.User,
	) (*model.User, error)
}

type GetUserQuery struct {
	ID      null.String
	AuthUID null.String
}
