package request_interceptor

import (
	"context"
	"time"
)

type contextKey string

const (
	requestTime contextKey = "requestTime"
)

func SetRequestTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, requestTime, t)
}

func GetRequestTime(ctx context.Context) time.Time {
	t, _ := ctx.Value(requestTime).(time.Time)
	return t
}
