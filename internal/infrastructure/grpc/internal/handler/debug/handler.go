package debug

import (
	"context"

	debug_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/debug_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase"
)

type DebugHandler struct {
	debugInteractor usecase.DebugInteractor
}

func NewDebugHandler(
	debugInteractor usecase.DebugInteractor,
) debug_apiv1.DebugV1ServiceServer {
	return &DebugHandler{
		debugInteractor: debugInteractor,
	}
}

func (h *DebugHandler) CreateStaffIDToken(ctx context.Context, req *debug_apiv1.CreateStaffIDTokenRequest) (*debug_apiv1.CreateStaffIDTokenResponse, error) {
	idToken, err := h.debugInteractor.CreateStaffIDToken(ctx, req.GetAuthUid(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &debug_apiv1.CreateStaffIDTokenResponse{
		IdToken: idToken,
	}, nil
}
