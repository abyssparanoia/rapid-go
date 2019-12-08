package service

import (
	"context"
	"net/http"

	"github.com/abyssparanoia/rapid-go/src/pkg/log"
)

func newUserNotExistError(ctx context.Context, userID string) error {
	return log.Errorc(ctx, http.StatusNotFound, "user ID %s does not exist", userID)
}
