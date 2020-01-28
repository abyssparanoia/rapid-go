package input

import (
	"net/http"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/abyssparanoia/rapid-go/push-notification/domain/model"

	"context"
)

// RegisterSetToken ... register set token input
type RegisterSetToken struct {
	AppID    string
	UserID   string
	Platform model.Platform
	DeviceID string
	Token    string
}

// NewRegisterSetToken ... new register set token input
func NewRegisterSetToken(
	ctx context.Context,
	appID string,
	userID string,
	platform string,
	deviceID string,
	token string) (*RegisterSetToken, error) {

	_platform, err := model.NewPlatform(ctx, platform)
	if err != nil {
		return nil, log.Errorc(ctx, http.StatusBadRequest, "model.NewPlatform")
	}

	return &RegisterSetToken{
		AppID:    appID,
		UserID:   userID,
		Platform: _platform,
		DeviceID: deviceID,
		Token:    token,
	}, nil
}
