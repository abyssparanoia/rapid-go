package mock_repository

import (
	context "context"

	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
)

type mockTx struct{}

func (mock *mockTx) ROTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (mock *mockTx) RWTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestMockTransactable() repository.Transactable {
	return &mockTx{}
}
