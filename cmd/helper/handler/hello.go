package handler

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
)

func (h *helperHandler) HelloWorld(ctx context.Context) {
	log.Debugf(ctx, "hello")
}
