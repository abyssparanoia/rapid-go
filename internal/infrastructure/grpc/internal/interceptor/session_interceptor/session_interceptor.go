package session_interceptor

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/usecase"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
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
	idToken, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil || idToken == "" {
		return ctx, nil
	}
	claims, err := i.authenticationInteractor.VerifyIDToken(ctx, input.NewVerifyIDToken(idToken))
	if err != nil {
		return ctx, err
	}
	ctx = SaveSessionContext(ctx, newSessionContext(claims))
	return ctx, nil
}
