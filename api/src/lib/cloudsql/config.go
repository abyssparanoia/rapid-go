package cloudsql

import (
	"fmt"
	"os"
	"strings"
)

// CSQLConfig ... CloudSQLの接続情報
type CSQLConfig struct {
	ConnectionName string
	User           string
	Password       string
}

// GetCSQLConfig ... CloudSQLの接続情報を取得する
func GetCSQLConfig(db string) *CSQLConfig {
	db = strings.ToUpper(db)

	cnKey := fmt.Sprintf("CLOUDSQL_%s_CONNECTION_NAME", db)
	cn := os.Getenv(cnKey)
	if cn == "" {
		panic(fmt.Errorf("no config key %s", cnKey))
	}

	uKey := fmt.Sprintf("CLOUDSQL_%s_USER", db)
	u := os.Getenv(uKey)
	if u == "" {
		panic(fmt.Errorf("no config key %s", uKey))
	}

	pKey := fmt.Sprintf("CLOUDSQL_%s_PASSWORD", db)
	p := os.Getenv(pKey)
	if p == "" {
		panic(fmt.Errorf("no config key %s", pKey))
	}

	return &CSQLConfig{
		ConnectionName: cn,
		User:           u,
		Password:       p,
	}
}
