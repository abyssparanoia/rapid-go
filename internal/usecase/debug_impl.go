package usecase

import (
	"context"

	"github.com/aarondl/null/v9"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
)

type debugInteractor struct {
	adminAuthenticationRepository repository.AdminAuthentication
	staffAuthenticationRepository repository.StaffAuthentication
}

func NewDebugInteractor(
	adminAuthenticationRepository repository.AdminAuthentication,
	staffAuthenticationRepository repository.StaffAuthentication,
) DebugInteractor {
	return &debugInteractor{
		adminAuthenticationRepository,
		staffAuthenticationRepository,
	}
}

func (i *debugInteractor) CreateAdminIDToken(
	ctx context.Context,
	authUID string,
	password string,
) (string, error) {
	return i.adminAuthenticationRepository.CreateIDToken(ctx, authUID, password)
}

func (i *debugInteractor) CreateStaffIDToken(
	ctx context.Context,
	authUID string,
	password string,
) (string, error) {
	return i.staffAuthenticationRepository.CreateIDToken(ctx, authUID, password)
}

func (i *debugInteractor) CreateStaffAuthUID(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	return i.staffAuthenticationRepository.CreateUser(ctx, repository.StaffAuthenticationCreateUserParam{
		Email:    email,
		Password: null.StringFrom(password),
	})
}
