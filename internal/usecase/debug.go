package usecase

import "context"

type DebugInteractor interface {
	CreateStaffIDToken(
		ctx context.Context,
		authUID string,
		password string,
	) (string, error)
}
