package usecase

import "context"

type DebugInteractor interface {
	CreateAdminIDToken(
		ctx context.Context,
		authUID string,
		password string,
	) (string, error)
	CreateStaffIDToken(
		ctx context.Context,
		authUID string,
		password string,
	) (string, error)
	CreateStaffAuthUID(
		ctx context.Context,
		email string,
		password string,
	) (string, error)
}
