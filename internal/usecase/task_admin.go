package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_usecase
type TaskAdminInteractor interface {
	Create(
		ctx context.Context,
		param *input.TaskCreateAdmin,
	) (*output.TaskCreateAdmin, error)
}
