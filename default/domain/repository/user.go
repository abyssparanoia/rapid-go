package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/default/domain/model"
)

// User ... user interface
type User interface {
	Get(ctx context.Context, userID string) (*model.User, error)
}
