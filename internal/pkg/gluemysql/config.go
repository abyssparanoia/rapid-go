package gluemysql

import (
	"fmt"
	"os"
)

// Config ... config for mysql
type Config struct {
	Host     string
	User     string
	Password string
	DB       string
}

// NewConfig ... get config for mysql from envelopment value
func NewConfig() *Config {

	hKey := fmt.Sprintf("DB_HOST")
	h := os.Getenv(hKey)
	if h == "" {
		panic(fmt.Errorf("no config key %s", hKey))
	}

	uKey := fmt.Sprintf("DB_USER")
	u := os.Getenv(uKey)
	if u == "" {
		panic(fmt.Errorf("no config key %s", uKey))
	}

	pKey := fmt.Sprintf("DB_PASSWORD")
	p := os.Getenv(pKey)
	if p == "" {
		panic(fmt.Errorf("no config key %s", pKey))
	}

	dKey := fmt.Sprintf("DB_DATABASE")
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
