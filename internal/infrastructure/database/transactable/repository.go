package transactable

import (
	"context"

	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/repository"
)

type transactable struct{}

func NewTransactable() repository.Transactable {
	return &transactable{}
}

func (r *transactable) RWTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return RunTx(ctx, fn)
}
