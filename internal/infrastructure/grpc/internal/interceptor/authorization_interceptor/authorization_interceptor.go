package authorization_interceptor

import (
	"context"
	"strings"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/session_interceptor"
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
		staffSessionContext, ok := session_interceptor.GetStaffSessionContext(ctx)
		method := info.FullMethod
		if strings.Contains(method, "AdminV1Service") {
			if !ok || !staffSessionContext.StaffRole.Value().IsAdmin() {
				return nil, errors.InvalidAdminRequestUserErr
			}
		}

		// TODO:
		// use casbin
		return handler(ctx, req)
	}
}
