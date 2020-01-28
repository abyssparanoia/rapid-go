package service

import (
	"github.com/abyssparanoia/rapid-go/push-notification/domain/model"

	"context"
)

// Token ... token service interface
type Token interface {
	Set(ctx context.Context,
		token *model.Token) error
}
