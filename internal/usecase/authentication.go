package usecase

import (
	"context"

	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/model"
	"github.com/playground-live/moala-meet-and-greet-back/internal/usecase/input"
)

type AuthenticationInteractor interface {
	VerifyIDToken(
		ctx context.Context,
		param *input.VerifyIDToken,
	) (*model.Claims, error)
}
