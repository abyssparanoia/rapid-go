package redis

import (
	"crypto/tls"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func NewClient(
	host string,
	port string,
	username string,
	password string,
	tlsEnable bool,
) *redis.Client {
	var tlsConfig *tls.Config = nil
	if tlsEnable {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	return redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%s", host, port),
		Username:   username,
		Password:   password,
		DB:         0,
		MaxRetries: 5,
		TLSConfig:  tlsConfig,
	})
}
