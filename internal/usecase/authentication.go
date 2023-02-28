package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type AuthenticationInteractor interface {
	VerifyStaffIDToken(
		ctx context.Context,
		param *input.VerifyIDToken,
	) (*model.StaffClaims, error)
}
