package handler

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/default-grpc/usecase"
	pb "github.com/abyssparanoia/rapid-go/proto/default"
)

// UserHandler ...
type UserHandler struct {
	userUsecase usecase.User
}

// GetUser ...
func (h *UserHandler) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (*pb.User, error) {
	user, err := h.userUsecase.Get(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:                  user.ID,
		DisplayName:         user.DisplayName,
		IconImagePath:       user.IconImagePath,
		BackgroundImagePath: user.BackgroundImagePath,
	}, nil
}

// NewUserHandler ...
func NewUserHandler(
	userUsecase usecase.User,
) *UserHandler {
	return &UserHandler{userUsecase}
}
