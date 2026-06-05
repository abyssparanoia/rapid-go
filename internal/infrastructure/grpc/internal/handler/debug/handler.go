package debug

import (
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
