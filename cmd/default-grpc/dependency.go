package main

import (
	"fmt"
	"net"

	"github.com/abyssparanoia/rapid-go/internal/default-grpc/handler"
	"github.com/abyssparanoia/rapid-go/internal/default-grpc/infrastructure/repository"
	"github.com/abyssparanoia/rapid-go/internal/default-grpc/usecase"
	"github.com/abyssparanoia/rapid-go/internal/pkg/gluemysql"
	grpc_requestlog "github.com/abyssparanoia/rapid-go/internal/pkg/interceptor/requestlog"
	pb "github.com/abyssparanoia/rapid-go/proto/default"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func recoveryFuncFactory(logger *zap.Logger) func(p interface{}) error {
	return func(p interface{}) error {
		logger.Error(fmt.Sprintf("p: %+v\n", p))
		return status.Errorf(codes.Internal, "Unexpected error: %+v\n", p)
	}
}

func newDefaultServer(logger *zap.Logger, e *environment) *grpc.Server {
	grpc_zap.ReplaceGrpcLogger(logger)

	_ = gluemysql.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase)

	// Repository
	userRepository := repository.NewUser()

	userUsecase := usecase.NewUser(userRepository)

	userHandler := handler.NewUserHandler(userUsecase)

	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryFuncFactory(logger)),
	}

	server := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_requestlog.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
	)

	pb.RegisterUserServiceServer(server, userHandler)
	reflection.Register(server)
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", e.Port))
	if err != nil {
		panic(err)
	}

	go func() {
		logger.Info(fmt.Sprintf("Listening grpc on %s:%s", "localhost", e.Port))
		if err := server.Serve(listen); err != nil {
			logger.Error("server.Serve", zap.Error(err))
		}
	}()

	return server
}
