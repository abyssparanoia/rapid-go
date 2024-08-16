package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/dependency"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	internal_grpc "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
	debug_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/debug_api/v1"
	public_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/public_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/http/internal/handler"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/http/internal/middlewares"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger"
	"github.com/caarlos0/env/v11"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/encoding/protojson"
)

const InternalGRPCPort = "50051"

type Utf8JsonMarshaller struct {
	runtime.JSONPb
}

func Run() {
	e := &environment.Environment{}
	if err := env.Parse(e); err != nil {
		panic(err)
	}

	logger := logger.New()

	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Second)
	defer cancel()

	addr := fmt.Sprintf(":%s", e.Port)

	d := &dependency.Dependency{}
	d.Inject(ctx, e)

	grpcServer := internal_grpc.NewServer(ctx, e, logger, d)
	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%s", InternalGRPCPort))
	if err != nil {
		panic(err)
	}

	// Run grpc server
	logger.Info(fmt.Sprintf("[START] server. port: %s\n", addr))

	go func() {
		if err = grpcServer.Serve(grpcLis); err != nil {
			logger.Error("failed to start server", zap.Error(err))
		}
	}()

	grpcGateway := runtime.NewServeMux(runtime.WithMarshalerOption("*", &runtime.HTTPBodyMarshaler{
		Marshaler: &CustomJSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	}))

	conn, err := grpc.NewClient(
		fmt.Sprintf(":%s", InternalGRPCPort),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                30 * time.Millisecond,
				Timeout:             20 * time.Millisecond,
				PermitWithoutStream: true,
			}),
	)
	if err != nil {
		panic(err)
	}

	if err = admin_apiv1.RegisterAdminV1ServiceHandler(context.Background(), grpcGateway, conn); err != nil {
		panic(err)
	}

	if err = public_apiv1.RegisterPublicV1ServiceHandler(context.Background(), grpcGateway, conn); err != nil {
		panic(err)
	}

	if e.Environment == "local" || e.Environment == "development" {
		if err = debug_apiv1.RegisterDebugV1ServiceHandler(context.Background(), grpcGateway, conn); err != nil {
			panic(err)
		}
	}

	if err = grpcGateway.HandlePath(http.MethodGet, "/", handler.Ping); err != nil {
		panic(err)
	}

	// server
	server := http.Server{
		Addr:              addr,
		Handler:           middlewares.CORS(grpcGateway),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err = server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("[CLOSED] server closed with error", zap.Error(err))
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	logger.Info(fmt.Sprintf("SIGNAL %d received, so server shutting down now...\n", <-quit))

	err = server.Shutdown(ctx)
	if err != nil {
		logger.Error("failed to gracefully shutdown", zap.Error(err))
	}
	grpcServer.GracefulStop()

	logger.Info("server shutdown completed")
}
