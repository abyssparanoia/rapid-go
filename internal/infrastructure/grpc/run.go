package grpc

import (
	"context"
	"fmt"
	debug_util "runtime/debug"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/dependency"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/handler/admin"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/handler/debug"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/handler/public"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/authorization_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/request_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/session_interceptor"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
	debug_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/debug_api/v1"
	public_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/public_api/v1"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func recoverFuncFactory(l *zap.Logger) func(p interface{}) error {
	return func(p interface{}) error {
		l.Error(fmt.Sprintf("p: %+v\n : %s", p, debug_util.Stack()))
		return status.Errorf(codes.Internal, "Unexpected error")
	}
}

func NewServer(ctx context.Context,
	e *environment.Environment,
	logger *zap.Logger,
	dependency *dependency.Dependency,
) *grpc.Server {

	requestLogInterceptor := request_interceptor.NewRequestLog(logger)
	authFunc := session_interceptor.NewSession(dependency.AuthenticationInteractor)
	authorizationInterceptor := authorization_interceptor.NewAuthorization()

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
				grpc_recovery.UnaryServerInterceptor(
					grpc_recovery.WithRecoveryHandler(recoverFuncFactory(logger)),
				),
				grpc_auth.UnaryServerInterceptor(authFunc.Authenticate),
				authorizationInterceptor.UnaryServerInterceptor(),
			),
		),
	)

	admin_apiv1.RegisterAdminV1ServiceServer(
		server,
		admin.NewAdminHandler(
			dependency.AdminTenantInteractor,
			dependency.AdminUserInteractor,
			dependency.AdminAssetInteractor,
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
