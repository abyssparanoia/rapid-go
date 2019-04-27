package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/src/infrastructure/entity"
)

// User ... user interface
type User interface {
	Get(ctx context.Context, userID int64) (*entity.User, error)
}
