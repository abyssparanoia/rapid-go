package session_interceptor

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger"
	"go.uber.org/zap"
)

type staffContextKey struct{}

var staffSessionContextKey staffContextKey = staffContextKey{}

func SaveStaffSessionContext(
	ctx context.Context,
	claims *model.StaffClaims,
) context.Context {
	logger.AddFields(ctx, zap.Any("staff.claims", claims))
	return context.WithValue(ctx, staffSessionContextKey, *claims)
}

func RequireStaffSessionContext(ctx context.Context) (*model.StaffClaims, error) {
	sctx, ok := GetStaffSessionContext(ctx)
	if !ok {
		return nil, errors.RequireStaffSessionErr.New()
	}
	return sctx, nil
}

func GetStaffSessionContext(ctx context.Context) (*model.StaffClaims, bool) {
	sessionContext, ok := ctx.Value(staffSessionContextKey).(model.StaffClaims)
	return &sessionContext, ok
}
