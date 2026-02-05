package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_usecase
type StaffMeInteractor interface {
	SignUp(
		ctx context.Context,
		param *input.StaffSignUp,
	) (*model.Staff, error)
	Get(
		ctx context.Context,
		param *input.StaffGetMe,
	) (*model.Staff, error)
	Update(
		ctx context.Context,
		param *input.StaffUpdateMe,
	) (*model.Staff, error)
}
