package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type adminAuthenticationInteractor struct {
	adminAuthenticationRepository repository.AdminAuthentication
}

func NewAdminAuthenticationInteractor(
	adminAuthenticationRepository repository.AdminAuthentication,
) AdminAuthenticationInteractor {
	return &adminAuthenticationInteractor{
		adminAuthenticationRepository,
	}
}

func (i *adminAuthenticationInteractor) VerifyAdminIDToken(
	ctx context.Context,
	param *input.VerifyIDToken,
) (*model.AdminClaims, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	return i.adminAuthenticationRepository.VerifyIDToken(ctx, param.IDToken)
}
