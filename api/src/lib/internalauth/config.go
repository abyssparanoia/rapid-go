package internalauth

import (
	"fmt"
	"os"
)

const (
	tokenKey = "INTERNAL_AUTH_TOKEN"
)

// GetToken ... 内部認証のTokenを取得する
func GetToken() string {
	k := os.Getenv(tokenKey)
	if k == "" {
		panic(fmt.Errorf("no token internal auth: %s", tokenKey))
	}
	return k
}
