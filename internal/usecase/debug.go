package usecase

import "context"

type DebugInteractor interface {
	CreateIDToken(
		ctx context.Context,
		authUID string,
		password string,
	) (string, error)
}
