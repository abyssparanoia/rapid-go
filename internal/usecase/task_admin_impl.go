package usecase

import (
	"context"

	"github.com/aarondl/null/v9"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

type taskAdminInteractor struct {
	transactable                  repository.Transactable
	adminRepository               repository.Admin
	adminAuthenticationRepository repository.AdminAuthentication
}

func NewTaskAdminInteractor(
	transactable repository.Transactable,
	adminRepository repository.Admin,
	adminAuthenticationRepository repository.AdminAuthentication,
) TaskAdminInteractor {
	return &taskAdminInteractor{
		transactable:                  transactable,
		adminRepository:               adminRepository,
		adminAuthenticationRepository: adminAuthenticationRepository,
	}
}

func (i *taskAdminInteractor) Create(
	ctx context.Context,
	param *input.TaskCreateAdmin,
) (*output.TaskCreateAdmin, error) {
	// 1. Validate input
	if err := param.Validate(); err != nil {
		return nil, err
	}

	var admin *model.Admin
	var authUID string

	// 2. Execute in transaction
	if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
		// 3. Create user in Cognito
		uid, err := i.adminAuthenticationRepository.CreateUser(
			ctx,
			repository.AdminAuthenticationCreateUserParam{
				Email:    param.Email,
				Password: null.StringFrom(param.Password),
			},
		)
		if err != nil {
			return err
		}
		authUID = uid

		// 4. Create Admin entity
		admin = model.NewAdmin(
			model.AdminRoleRoot,
			authUID,
			param.Email,
			param.DisplayName,
			param.RequestTime,
		)

		// 5. Save to DB
		if err := i.adminRepository.Create(ctx, admin); err != nil {
			return err
		}

		// 6. Store claims in Cognito
		claims := model.NewAdminClaims(
			authUID,
			param.Email,
			null.StringFrom(admin.ID),
			nullable.TypeFrom(admin.Role),
		)
		if err := i.adminAuthenticationRepository.StoreClaims(ctx, authUID, claims); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// 7. Return result
	return output.NewTaskCreateAdmin(admin.ID, authUID, param.Password), nil
}
