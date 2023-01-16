package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/volatiletech/null/v8"
)

type adminUserInteractor struct {
	transactable     repository.Transactable
	tenantRepository repository.Tenant
	userService      service.User
}

func NewAdminUserInteractor(
	transactable repository.Transactable,
	tenantRepository repository.Tenant,
	userService service.User,
) AdminUserInteractor {
	return &adminUserInteractor{
		transactable,
		tenantRepository,
		userService,
	}
}

func (i *adminUserInteractor) Create(
	ctx context.Context,
	param *input.AdminCreateUser,
) (*model.User, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}

	var user *model.User
	if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
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
			return err
		}

		user, err = i.userService.Create(
			ctx,
			service.UserCreateParam{
				TenantID:    tenant.ID,
				Email:       param.Email,
				Password:    "random1234",
				UserRole:    param.Role,
				DisplayName: param.DisplayName,
				ImagePath:   "user_profile_images/default_image.jpeg",
				RequestTime: param.RequestTime,
			},
		)
		if err != nil {
			return err
		}

		user.Tenant = tenant

		return nil
	}); err != nil {
		return nil, err
	}

	return user, nil
}
