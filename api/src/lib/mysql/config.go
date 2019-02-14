package mysql

import (
	"fmt"
	"os"
)

// CSQLConfig ... MySQLの接続情報
type SQLConfig struct {
	ConnectionName string
	User           string
	Password       string
	Database       string
}

// GetSQLConfig ... MySQLの接続情報を取得する
func GetSQLConfig() *SQLConfig {
	cnKey := "DB_HOST"
	cn := os.Getenv(cnKey)
	if cn == "" {
		panic(fmt.Errorf("no config key %s", cnKey))
	}

	uKey := "DB_USER"
	u := os.Getenv(uKey)
	if u == "" {
		panic(fmt.Errorf("no config key %s", uKey))
	}

	pKey := "DB_PASSWORD"
	p := os.Getenv(pKey)
	if p == "" {
		panic(fmt.Errorf("no config key %s", pKey))
	}

	dKey := "DB_DATABASE"
	d := os.Getenv(dKey)
	if d == "" {
		panic(fmt.Errorf("no config key %s", dKey))
	}

	return &SQLConfig{
		ConnectionName: cn,
		User:           u,
		Password:       p,
		Database:       d,
	}
}
