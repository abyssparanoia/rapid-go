package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type authenticationInteractor struct {
	authenticationRepository repository.Authentication
}

func NewAuthenticationInteractor(
	authenticationRepository repository.Authentication,
) AuthenticationInteractor {
	return &authenticationInteractor{
		authenticationRepository,
	}
}

func (i *authenticationInteractor) VerifyIDToken(
	ctx context.Context,
	param *input.VerifyIDToken,
) (*model.Claims, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	return i.authenticationRepository.VerifyIDToken(ctx, param.IDToken)
}
