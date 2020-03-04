package input

import (
	"net/http"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/domain/model"

	"context"
)

// TokenSet ... register set token input
type TokenSet struct {
	AppID    string
	UserID   string
	Platform model.Platform
	DeviceID string
	Token    string
}

// NewTokenSet ... new register set token input
func NewTokenSet(
	ctx context.Context,
	appID string,
	userID string,
	platform string,
	deviceID string,
	token string) (*TokenSet, error) {

	_platform, err := model.NewPlatform(ctx, platform)
	if err != nil {
		return nil, log.Errorc(ctx, http.StatusBadRequest, "model.NewPlatform")
	}

	return &TokenSet{
		AppID:    appID,
		UserID:   userID,
		Platform: _platform,
		DeviceID: deviceID,
		Token:    token,
	}, nil
}

// TokenDelete ... register delete token input
type TokenDelete struct {
	AppID    string
	UserID   string
	Platform model.Platform
	DeviceID string
}

// NewTokenDelete ... new register delete token input
func NewTokenDelete(
	ctx context.Context,
	appID string,
	userID string,
	platform string,
	deviceID string,
) (*TokenDelete, error) {

	_platform, err := model.NewPlatform(ctx, platform)
	if err != nil {
		return nil, log.Errorc(ctx, http.StatusBadRequest, "model.NewPlatform")
	}

	return &TokenDelete{
		AppID:    appID,
		UserID:   userID,
		Platform: _platform,
		DeviceID: deviceID,
	}, nil
}
