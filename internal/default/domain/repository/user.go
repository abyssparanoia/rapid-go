package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/default/domain/model"
)

// User ... user interface
type User interface {
	Get(
		ctx context.Context,
		userID string,
	) (*model.User, error)
	Create(
		ctx context.Context,
		user *model.User,
	) (*model.User, error)
}
