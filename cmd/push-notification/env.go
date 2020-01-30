package main

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Environment ... environment variable
type Environment struct {
	Port       string `envconfig:"PORT"                           default:"8080"`
	ENV        string `envconfig:"ENV"                         required:"true"`
	ProjectID  string `envconfig:"PROJECT_ID"                     required:"true"`
	LocationID string `envconfig:"LOCATION_ID"                    default:"asia-northeast1"`
	// ServiceID       string `envconfig:"SERVICE_ID"                     required:"true"`
	CredentialsPath string `envconfig:"GOOGLE_APPLICATION_CREDENTIALS" required:"true"`
	MinLogSeverity  string `envconfig:"MIN_LOG_SEVERITY"               required:"true"`
	FcmServerKey    string `envconfig:"FCM_SERVER_KEY" required:"true"`
}

// Get ... get env
func (e *Environment) Get() {
	err := godotenv.Load(".env.push-notification")
	if err != nil {
		panic(err)
	}
	err = envconfig.Process("", e)
	if err != nil {
		panic(err)
	}
}
