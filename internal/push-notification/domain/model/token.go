package model

import (
	"context"

	"github.com/pkg/errors"

	"github.com/abyssparanoia/rapid-go/internal/pkg/util"
)

// Token ... token model
type Token struct {
	ID        string
	Platform  Platform
	AppID     string
	UserID    string
	DeviceID  string
	Value     string
	CreatedAt int64
}

// Exists ... check exists or not
func (m *Token) Exists() bool {
	return m != nil
}

// NewToken ... new token model
func NewToken(platform Platform,
	appID, deviceID, value string,
) *Token {
	return &Token{
		Platform:  platform,
		AppID:     appID,
		DeviceID:  deviceID,
		Value:     value,
		CreatedAt: util.TimeNowUnixMill(),
	}
}

// NewTokenValues ... get token value list from token model list
func NewTokenValues(tokens []*Token) []string {
	tokenValues := make([]string, len(tokens))
	for index, token := range tokens {
		tokenValues[index] = token.Value
	}
	return tokenValues
}

// Platform ... platform type
type Platform string

const (
	// PlatformIOS ... ios
	PlatformIOS Platform = "ios"
	// PlatformAndroid ... android
	PlatformAndroid Platform = "android"
	// PlatformWeb ... web
	PlatformWeb Platform = "web"
)

// NewPlatform ... new platform
func NewPlatform(ctx context.Context, platform string) (Platform, error) {

	switch Platform(platform) {
	case PlatformIOS, PlatformAndroid, PlatformWeb:
		return Platform(platform), nil
	default:
	}

	return "", errors.New("no match case")
}

// MustPlatform ... must platform
func MustPlatform(platform string) Platform {

	switch Platform(platform) {
	case PlatformIOS, PlatformAndroid, PlatformWeb:
		return Platform(platform)
	default:
	}

	panic("no match case")
}

func (v Platform) String() string {
	return string(v)
}
