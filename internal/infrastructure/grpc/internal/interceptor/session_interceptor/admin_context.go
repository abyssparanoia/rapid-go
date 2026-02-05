package session_interceptor

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
)

type adminContextKey struct{}

var adminSessionContextKey adminContextKey = adminContextKey{}

func SaveAdminSessionContext(ctx context.Context, claims *model.AdminClaims) context.Context {
	return context.WithValue(ctx, adminSessionContextKey, claims)
}

func RequireAdminSessionContext(ctx context.Context) (*model.AdminClaims, error) {
	claims, ok := GetAdminSessionContext(ctx)
	if !ok {
		return nil, errors.RequireAdminSessionErr.New()
	}
	return claims, nil
}

func GetAdminSessionContext(ctx context.Context) (*model.AdminClaims, bool) {
	claims, ok := ctx.Value(adminSessionContextKey).(*model.AdminClaims)
	return claims, ok
}
