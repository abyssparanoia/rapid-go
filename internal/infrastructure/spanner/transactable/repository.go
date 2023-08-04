package transactable

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/spanner/internal/dbmodel"
)

type transactable struct {
	*dbmodel.SpannerTransactable
}

func NewTransactable(
	spannerCli *spanner.Client,
) repository.Transactable {
	t := dbmodel.NewTransactable(spannerCli)
	return &transactable{
		t,
	}
}

func (r *transactable) ROTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if err := r.SpannerTransactable.ROTx(ctx, fn); err != nil {
		return err
	}
	return nil
}

func (r *transactable) RWTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if err := r.SpannerTransactable.RWTx(ctx, fn); err != nil {
		return err
	}
	return nil
}
