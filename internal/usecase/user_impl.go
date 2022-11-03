package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/volatiletech/null/v8"
)

type userInteractor struct {
	transactable             repository.Transactable
	userRepository           repository.User
	tenantRepository         repository.Tenant
	authenticationRepository repository.Authentication
}

func NewUserInteractor(
	transactable repository.Transactable,
	userRepository repository.User,
	tenantRepository repository.Tenant,
	authenticationRepository repository.Authentication,
) UserInteractor {
	return &userInteractor{
		transactable,
		userRepository,
		tenantRepository,
		authenticationRepository,
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

		res, err := i.authenticationRepository.GetUserByEmail(ctx, param.Email)
		if err != nil {
			return err
		}
		var authUID string
		// 存在してない場合、新規作成する
		if !res.Exist {
			authUID, err = i.authenticationRepository.CreateUser(
				ctx,
				repository.AuthenticationCreateUserParam{
					Email:    param.Email,
					Password: null.StringFrom(param.Passoword),
				},
			)
			if err != nil {
				return err
			}
		} else {
			authUID = res.AuthUID
		}

		user := model.NewUser(
			tenant.ID,
			model.UserRoleAdmin,
			authUID,
			"Root User",
			"user_profile_images/default_image.jpeg",
			param.Email,
			param.RequestTime,
		)

		if _, err := i.userRepository.Create(ctx, user); err != nil {
			return err
		}

		claims := model.NewClaims(authUID)
		claims.SetTenantID(tenant.ID)
		claims.SetUserID(user.ID)
		claims.SetUserRole(user.Role)
		if err := i.authenticationRepository.StoreClaims(ctx, authUID, claims); err != nil {
			return err
		}

		return nil
	})
}
