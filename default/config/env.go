package config

import (
	"os"
)

// IsEnvDeveloping ... development
func IsEnvDeveloping() bool {
	return os.Getenv("ENV") == "develop"
}

// IsEnvStaging ... staging
func IsEnvStaging() bool {
	return os.Getenv("ENV") == "staging"
}

// IsEnvProduction ... production
func IsEnvProduction() bool {
	return os.Getenv("ENV") == "production"
}
