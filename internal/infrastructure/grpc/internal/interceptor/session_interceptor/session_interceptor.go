package session_interceptor

import (
	"context"
	"strings"

	"github.com/abyssparanoia/rapid-go/internal/usecase"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
)

type Session struct {
	authenticationInteractor usecase.AuthenticationInteractor
}

func NewSession(
	authenticationInteractor usecase.AuthenticationInteractor,
) *Session {
	return &Session{
		authenticationInteractor,
	}
}

func (i *Session) Authenticate(ctx context.Context) (context.Context, error) {
	method, _ := grpc.Method(ctx)
	idToken, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil || idToken == "" {
		return ctx, nil //nolint:nilerr // ignore
	}
	if strings.Contains(method, "AdminV1Service") {
		claims, err := i.authenticationInteractor.VerifyStaffIDToken(ctx, input.NewVerifyIDToken(idToken))
		if err != nil {
			return ctx, err
		}
		ctx = SaveStaffSessionContext(ctx, claims)
	}
	return ctx, nil
}
