package basicauth

import (
	"fmt"
	"os"
)

const (
	userKey     = "BASIC_AUTH_USER"
	passwordKey = "BASIC_AUTH_PASSWORD"
)

// Account ... ベーシック認証のアカウント
type Account struct {
	User     string
	Password string
}

// GetAccount ... ベーシック認証のアカウントを取得する
func GetAccount() *Account {
	u := os.Getenv(userKey)
	if u == "" {
		panic(fmt.Errorf("no account basic auth user: %s", userKey))
	}
	p := os.Getenv(passwordKey)
	if p == "" {
		panic(fmt.Errorf("no account basic auth password: %s", passwordKey))
	}
	return &Account{
		User:     u,
		Password: p,
	}
}
