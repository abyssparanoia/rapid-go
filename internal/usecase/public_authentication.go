package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type PublicAuthenticationInteractor interface {
	SignIn(
		ctx context.Context,
		param *input.PublicSignIn,
	) (*model.User, error)
}
