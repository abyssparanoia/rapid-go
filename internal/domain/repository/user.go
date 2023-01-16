package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/volatiletech/null/v8"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type User interface {
	Get(
		ctx context.Context,
		query GetUserQuery,
	) (*model.User, error)
	Create(
		ctx context.Context,
		user *model.User,
	) (*model.User, error)
}

type GetUserQuery struct {
	BaseGetOptions
	ID      null.String
	AuthUID null.String
}
