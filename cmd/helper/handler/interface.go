package handler

import "context"

// HelperHandler ...
type HelperHandler interface {
	HelloWorld(ctx context.Context)
}

type helperHandler struct{}

// NewHelperHandler ... new helper handler
func NewHelperHandler() HelperHandler {
	return &helperHandler{}
}
