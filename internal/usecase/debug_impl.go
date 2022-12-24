package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
)

type debugInteractor struct {
	authenticationRepository repository.Authentication
}

func NewDebugInteractor(
	authenticationRepository repository.Authentication,
) DebugInteractor {
	return &debugInteractor{
		authenticationRepository,
	}
}

func (i *debugInteractor) CreateIDToken(ctx context.Context, authUID string) (string, error) {
	return i.authenticationRepository.CreateIDToken(ctx, authUID)
}
