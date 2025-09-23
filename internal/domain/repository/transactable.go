package repository

import "context"

//go:generate go tool go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type Transactable interface {
	ROTx(ctx context.Context, fn func(ctx context.Context) error) error
	RWTx(ctx context.Context, fn func(ctx context.Context) error) error
}
