package mysql

import (
	"fmt"
	"os"
	"strings"
)

// Config ... config for mysql
type Config struct {
	Host     string
	User     string
	Password string
	DB       string
}

// NewConfig ... get config for mysql from envelopment value
func NewConfig(db string) *Config {
	db = strings.ToUpper(db)

	hKey := fmt.Sprintf("MYSQL_%s_HOST", db)
	h := os.Getenv(hKey)
	if h == "" {
		panic(fmt.Errorf("no config key %s", hKey))
	}

	uKey := fmt.Sprintf("MYSQL_%s_USER", db)
	u := os.Getenv(uKey)
	if u == "" {
		panic(fmt.Errorf("no config key %s", uKey))
	}

	pKey := fmt.Sprintf("MYSQL_%s_PASSWORD", db)
	p := os.Getenv(pKey)
	if p == "" {
		panic(fmt.Errorf("no config key %s", pKey))
	}

	dKey := fmt.Sprintf("MYSQL_%s_DATABASE", db)
	d := os.Getenv(dKey)
	if d == "" {
		panic(fmt.Errorf("no config key %s", dKey))
	}

	return &Config{
		Host:     h,
		User:     u,
		Password: p,
		DB:       d,
	}
}
