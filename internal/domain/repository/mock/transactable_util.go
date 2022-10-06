package mock_repository

import (
	context "context"

	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/repository"
)

type mockTx struct{}

func (mock *mockTx) RWTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestMockTransactable() repository.Transactable {
	return &mockTx{}
}
