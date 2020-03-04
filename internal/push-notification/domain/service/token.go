package service

import (
	"github.com/abyssparanoia/rapid-go/internal/push-notification/domain/model"

	"context"
)

// Token ... token service interface
type Token interface {
	Set(ctx context.Context,
		token *model.Token) error
}
