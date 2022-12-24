package debug

import (
	"context"

	debug_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/debug_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase"
)

type DebugHander struct {
	debug_apiv1.UnimplementedDebugV1ServiceServer
	debugInteractor usecase.DebugInteractor
}

func NewDebugHandler(
	debugInteractor usecase.DebugInteractor,
) debug_apiv1.DebugV1ServiceServer {
	return &DebugHander{
		debugInteractor: debugInteractor,
	}
}

func (h *DebugHander) DebugCreateIDToken(ctx context.Context, req *debug_apiv1.DebugCreateIDTokenRequest) (*debug_apiv1.DebugCreateIDTokenResponse, error) {
	idToken, err := h.debugInteractor.CreateIDToken(ctx, req.GetAuthUid(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &debug_apiv1.DebugCreateIDTokenResponse{
		IdToken: idToken,
	}, nil
}
