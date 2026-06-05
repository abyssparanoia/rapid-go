package usecase

import "context"

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_usecase
type DebugInteractor interface {
	CreateAdminIDToken(
		ctx context.Context,
		email string,
		password string,
	) (string, error)
	CreateStaffIDToken(
		ctx context.Context,
		email string,
		password string,
	) (string, error)
	CreateStaffAuthUID(
		ctx context.Context,
		email string,
		password string,
	) (string, error)
}
