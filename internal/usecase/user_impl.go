package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type userInteractor struct {
	transactable     repository.Transactable
	tenantRepository repository.Tenant
	userService      service.User
}

func NewUserInteractor(
	transactable repository.Transactable,
	tenantRepository repository.Tenant,
	userService service.User,
) UserInteractor {
	return &userInteractor{
		transactable,
		tenantRepository,
		userService,
	}
}

func (i *userInteractor) CreateRoot(
	ctx context.Context,
	param *input.CreateRootUser,
) error {
	return i.transactable.RWTx(ctx, func(ctx context.Context) error {
		if err := param.Validate(); err != nil {
			return err
		}

		tenant := model.NewTenant("Platformer", param.RequestTime)
		if _, err := i.tenantRepository.Create(ctx, tenant); err != nil {
			return err
		}

		if _, err := i.userService.Create(
			ctx,
			service.UserCreateParam{
				TenantID:    tenant.ID,
				Email:       param.Email,
				Password:    "random1234",
				UserRole:    model.UserRoleAdmin,
				DisplayName: "Root User",
				ImagePath:   "user_profile_images/default_image.jpeg",
				RequestTime: param.RequestTime,
			},
		); err != nil {
			return err
		}

		return nil
	})
}
