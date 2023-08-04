package transactable

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
)

type transactable struct{}

func NewTransactable() repository.Transactable {
	return &transactable{}
}

func (r *transactable) ROTx(ctx context.Context, fn func(ctx context.Context) error) error {
	panic("not implement for mysql")
}

func (r *transactable) RWTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return RunTx(ctx, fn)
}
