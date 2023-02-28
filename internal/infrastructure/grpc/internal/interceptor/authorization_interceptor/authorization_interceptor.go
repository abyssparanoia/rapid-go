package authorization_interceptor

import (
	"context"
	"strings"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/session_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/pkg/errors"
	"google.golang.org/grpc"
)

type Authorization struct{}

func NewAuthorization() *Authorization {
	return &Authorization{}
}

func (i *Authorization) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (
		interface{},
		error,
	) {
		sessionContext, ok := session_interceptor.GetSessionContext(ctx)
		method := info.FullMethod
		if strings.Contains(method, "AdminV1Service") {
			if !ok || !sessionContext.StaffRole.Value.IsAdmin() {
				return nil, errors.UnauthorizedErr.Errorf("Invalid request staff")
			}
		}

		// TODO:
		// use casbin
		return handler(ctx, req)
	}
}
