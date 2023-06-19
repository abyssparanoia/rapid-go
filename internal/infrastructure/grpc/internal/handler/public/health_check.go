package public

import (
	"context"

	public_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/public_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger"
	"go.uber.org/zap"
)

func (h *PublicHandler) PublicDeepHealthCheck(ctx context.Context, req *public_apiv1.PublicDeepHealthCheckRequest) (*public_apiv1.PublicDeepHealthCheckResponse, error) {
	databaseStatus := "up"
	if err := h.databaseCli.DB.Ping(); err != nil {
		logger.L(ctx).Error("failed to  h.databaseCli.DB.Ping", zap.Error(err))
		databaseStatus = "down"
	}

	return &public_apiv1.PublicDeepHealthCheckResponse{
		DatabaseStatus: databaseStatus,
	}, nil
}
