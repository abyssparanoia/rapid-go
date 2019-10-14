package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/src/domain/model"
)

// User ... user interface
type User interface {
	Get(ctx context.Context, userID string) (*model.User, error)
}
