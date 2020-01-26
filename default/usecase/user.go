package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/default/domain/model"
)

// User ... inteface of User usecase
type User interface {
	Get(ctx context.Context, userID string) (*model.User, error)
}
