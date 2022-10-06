package usecase

import (
	"context"

	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/model"
	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/repository"
	"github.com/playground-live/moala-meet-and-greet-back/internal/usecase/input"
)

type adminUserInteractor struct {
	transactable             repository.Transactable
	authenticationRepository repository.Authentication
	userRepository           repository.User
	tenantRepository         repository.Tenant
}

func NewAdminUserInteractor(
	transactable repository.Transactable,
	authenticationRepository repository.Authentication,
	userRepository repository.User,
	tenantRepository repository.Tenant,
) AdminUserInteractor {
	return &adminUserInteractor{
		transactable,
		authenticationRepository,
		userRepository,
		tenantRepository,
	}
}

func (i *adminUserInteractor) Create(
	ctx context.Context,
	param *input.AdminCreateUser,
) (*model.User, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	tenant, err := i.tenantRepository.Get(
		ctx,
		param.TenantID,
		true,
	)
	if err != nil {
		return nil, err
	}

	authUID, err := i.authenticationRepository.CreateUser(
		ctx,
		repository.AuthenticationCreateUserParam{
			Email: param.Email,
		},
	)
	if err != nil {
		return nil, err
	}

	user := model.NewUser(
		param.TenantID,
		param.Role,
		authUID,
		param.DisplayName,
		"user_profile_images/default_image.jpeg",
		param.Email,
		param.RequestTime,
	)

	user, err = i.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	claims := model.NewClaims(authUID)
	claims.SetTenantID(param.TenantID)
	claims.SetUserID(user.ID)
	claims.SetUserRole(user.Role)
	if err := i.authenticationRepository.StoreClaims(ctx, authUID, claims); err != nil {
		return nil, err
	}

	user.Tenant = tenant

	return user, nil
}
