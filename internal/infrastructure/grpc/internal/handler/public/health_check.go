package public

import (
	"context"

	public_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/public_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger/logger_field"
)

func (h *PublicHandler) DeepHealthCheck(ctx context.Context, req *public_apiv1.DeepHealthCheckRequest) (*public_apiv1.DeepHealthCheckResponse, error) {
	databaseStatus := "up"
	if err := h.databaseCli.DB.Ping(); err != nil {
		logger.L(ctx).Error("failed to  h.databaseCli.DB.Ping", logger_field.Error(err))
		databaseStatus = "down"
	}

	return &public_apiv1.DeepHealthCheckResponse{
		DatabaseStatus: databaseStatus,
	}, nil
}
