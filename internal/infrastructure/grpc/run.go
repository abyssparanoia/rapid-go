package grpc

import (
	"context"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/dependency"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/handler/admin"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/handler/debug"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/handler/public"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/request_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/session_interceptor"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/mmg/admin_api/v1"
	debug_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/mmg/debug_api/v1"
	public_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/mmg/public_api/v1"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpcrecovry "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

func NewServer(ctx context.Context,
	e *environment.Environment,
	logger *zap.Logger,
	dependency *dependency.Dependency,
) *grpc.Server {

	requestLogInterceptor := request_interceptor.NewRequestLog(logger)
	authFunc := session_interceptor.NewSession(dependency.AuthenticationInteractor)

	server := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             10 * time.Second,
				PermitWithoutStream: true,
			},
		),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionIdle:     0,
				MaxConnectionAge:      10 * time.Minute,
				MaxConnectionAgeGrace: 0,
				Time:                  20 * time.Second,
				Timeout:               10 * time.Second,
			},
		),
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				requestLogInterceptor.UnaryServerInterceptor(),
				grpcrecovry.UnaryServerInterceptor(),
				grpc_auth.UnaryServerInterceptor(authFunc.Authenticate),
			),
		),
	)

	admin_apiv1.RegisterAdminV1ServiceServer(
		server,
		admin.NewAdminHandler(
			dependency.AdminTenantInteractor,
			dependency.AdminUserInteractor,
		),
	)
	public_apiv1.RegisterPublicV1ServiceServer(
		server,
		public.NewPublicHandler(
			dependency.PublicAuthenticationInteractor,
			dependency.PublicTenantInteractor,
		),
	)

	if e.Environment == "local" || e.Environment == "development" {
		debug_apiv1.RegisterDebugV1ServiceServer(
			server,
			debug.NewDebugHandler(
				dependency.FirebaseClient,
				e.FirebaseClientAPIKey,
			),
		)
	}

	reflection.Register(server)

	return server
}
