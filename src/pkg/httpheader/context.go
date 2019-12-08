package httpheader

import (
	"context"
)

type contextKey string

const paramsContextKey contextKey = "httpheader:params"

// GetParams ... get parameter from http header
func GetParams(ctx context.Context) Params {
	return ctx.Value(paramsContextKey).(Params)
}

func setParams(ctx context.Context, params Params) context.Context {
	return context.WithValue(ctx, paramsContextKey, params)
}
