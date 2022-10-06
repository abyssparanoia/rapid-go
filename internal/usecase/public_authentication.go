package usecase

import (
	"context"

	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/model"
	"github.com/playground-live/moala-meet-and-greet-back/internal/usecase/input"
)

type PublicAuthenticationInteractor interface {
	SignIn(
		ctx context.Context,
		param *input.PublicSignIn,
	) (*model.User, error)
}
