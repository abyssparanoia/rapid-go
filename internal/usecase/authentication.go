package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type AuthenticationInteractor interface {
	VerifyIDToken(
		ctx context.Context,
		param *input.VerifyIDToken,
	) (*model.Claims, error)
}
