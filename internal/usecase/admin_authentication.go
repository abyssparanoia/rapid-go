package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_usecase
type AdminAuthenticationInteractor interface {
	VerifyAdminIDToken(
		ctx context.Context,
		param *input.VerifyIDToken,
	) (*model.AdminClaims, error)
}
