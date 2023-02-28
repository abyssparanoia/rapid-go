package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type staffAuthenticationInteractor struct {
	staffAuthenticationRepository repository.StaffAuthentication
}

func NewAuthenticationInteractor(
	staffAuthenticationRepository repository.StaffAuthentication,
) AuthenticationInteractor {
	return &staffAuthenticationInteractor{
		staffAuthenticationRepository,
	}
}

func (i *staffAuthenticationInteractor) VerifyStaffIDToken(
	ctx context.Context,
	param *input.VerifyIDToken,
) (*model.StaffClaims, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	return i.staffAuthenticationRepository.VerifyIDToken(ctx, param.IDToken)
}
