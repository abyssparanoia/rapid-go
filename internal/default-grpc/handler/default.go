package handler

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/default-grpc/usecase"
	pb "github.com/abyssparanoia/rapid-go/proto/default"
)

// DefaultHandler ...
type DefaultHandler struct {
	userUsecase usecase.User
}

// GetUser ...
func (h *DefaultHandler) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (*pb.User, error) {
	return &pb.User{
		Id: "user_id",
	}, nil
}

// NewDefaultHandler ...
func NewDefaultHandler(
	userUsecase usecase.User,
) *DefaultHandler {
	return &DefaultHandler{userUsecase}
}
