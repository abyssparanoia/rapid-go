package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/volatiletech/null/v8"
)

type publicAuthenticationInteractor struct {
	userRepository repository.User
}

func NewPublicAuthenticationInteractor(
	userRepository repository.User,
) PublicAuthenticationInteractor {
	return &publicAuthenticationInteractor{
		userRepository,
	}
}

func (i *publicAuthenticationInteractor) SignIn(
	ctx context.Context,
	param *input.PublicSignIn,
) (*model.User, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	user, err := i.userRepository.Get(
		ctx,
		repository.GetUserQuery{
			AuthUID: null.StringFrom(param.AuthUID),
		},
		true,
		true,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
