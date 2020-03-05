package ctxhelper

import "context"

// Context ...
var Context context.Context

// SetContext ...
func SetContext(ctx context.Context) {
	Context = ctx
}

// GetContext ...
func GetContext() context.Context {
	return Context
}
