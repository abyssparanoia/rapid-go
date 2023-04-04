package session_interceptor

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/volatiletech/null/v8"
)

type StaffSessionContext struct {
	AuthUID   string
	TenantID  null.String
	StaffID   null.String
	StaffRole nullable.Type[model.StaffRole]
}

type staffContextKey struct{}

func newStaffSessionContext(claims *model.StaffClaims) *StaffSessionContext {
	return &StaffSessionContext{
		AuthUID:   claims.AuthUID,
		TenantID:  claims.TenantID,
		StaffID:   claims.StaffID,
		StaffRole: claims.StaffRole,
	}
}

var (
	staffSessionContextKey staffContextKey = staffContextKey{}
)

func SaveStaffSessionContext(
	ctx context.Context,
	sessionContext *StaffSessionContext,
) context.Context {
	return context.WithValue(ctx, staffSessionContextKey, *sessionContext)
}

func RequireStaffSessionContext(ctx context.Context) (*StaffSessionContext, error) {
	sctx, ok := GetStaffSessionContext(ctx)
	if !ok {
		return nil, errors.RequireStaffSessionErr.New()
	}
	return sctx, nil
}

func GetStaffSessionContext(ctx context.Context) (*StaffSessionContext, bool) {
	sessionContext, ok := ctx.Value(staffSessionContextKey).(StaffSessionContext)
	return &sessionContext, ok
}
