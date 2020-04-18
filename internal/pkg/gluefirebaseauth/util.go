package gluefirebaseauth

import (
	"context"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

const (
	headerPrefix      string = "Bearer"
	debugHeaderPrefix string = "user="
)

func getAuthClient(ctx context.Context) (*auth.Client, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}
	c, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func getTokenByAuthHeader(ah string) string {
	pLen := len(headerPrefix)
	if len(ah) > pLen && strings.ToUpper(ah[0:pLen]) == headerPrefix {
		return ah[pLen+1:]
	}
	return ""
}

func getUserByAuthHeader(ah string) string {
	if strings.HasPrefix(ah, debugHeaderPrefix) {
		return ah[len(debugHeaderPrefix):]
	}
	return ""
}
