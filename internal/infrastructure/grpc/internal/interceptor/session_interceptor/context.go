package session_interceptor

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/errors"
	"github.com/volatiletech/null/v8"
)

type SessionContext struct {
	AuthUID  string
	TenantID null.String
	UserID   null.String
	UserRole model.UserRole
}

type contextKey struct{}

func newSessionContext(claims *model.Claims) *SessionContext {
	return &SessionContext{
		AuthUID:  claims.AuthUID,
		TenantID: claims.TenantID,
		UserID:   claims.UserID,
		UserRole: claims.UserRole,
	}
}

var (
	sessionContextKey contextKey = contextKey{}
)

func SaveSessionContext(
	ctx context.Context,
	sessionContext *SessionContext,
) context.Context {
	return context.WithValue(ctx, sessionContextKey, *sessionContext)
}

func RequireSessionContext(ctx context.Context) (*SessionContext, error) {
	sctx, ok := GetSessionContext(ctx)
	if !ok {
		return nil, errors.UnauthorizedErr.New()
	}
	return sctx, nil
}

func GetSessionContext(ctx context.Context) (*SessionContext, bool) {
	sessionContext, ok := ctx.Value(sessionContextKey).(SessionContext)
	return &sessionContext, ok
}
