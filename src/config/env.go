package config

import (
	"os"
)

// IsEnvDeveloping ... 現在の環境がdevelopか判定する
func IsEnvDeveloping() bool {
	return os.Getenv("ENV") == "develop"
}

// IsEnvStaging ... 現在の環境がステージングか判定する
func IsEnvStaging() bool {
	return os.Getenv("ENV") == "staging"
}

// IsEnvProduction ... 現在の環境が本番か判定する
func IsEnvProduction() bool {
	return os.Getenv("ENV") == "production"
}
