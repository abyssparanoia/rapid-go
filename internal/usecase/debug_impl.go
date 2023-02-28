package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
)

type debugInteractor struct {
	staffAuthenticationRepository repository.StaffAuthentication
}

func NewDebugInteractor(
	staffAuthenticationRepository repository.StaffAuthentication,
) DebugInteractor {
	return &debugInteractor{
		staffAuthenticationRepository,
	}
}

func (i *debugInteractor) CreateStaffIDToken(
	ctx context.Context,
	authUID string,
	password string,
) (string, error) {
	return i.staffAuthenticationRepository.CreateIDToken(ctx, authUID, password)
}
