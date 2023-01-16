package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/volatiletech/null/v8"
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
		repository.GetTenantQuery{
			ID: null.StringFrom(param.TenantID),
			BaseGetOptions: repository.BaseGetOptions{
				OrFail: true,
			},
		},
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

	claims := model.NewClaims(
		authUID,
		null.StringFrom(param.TenantID),
		null.StringFrom(user.ID),
		nullable.TypeFrom(user.Role),
	)
	if err := i.authenticationRepository.StoreClaims(ctx, authUID, claims); err != nil {
		return nil, err
	}

	user.Tenant = tenant

	return user, nil
}
