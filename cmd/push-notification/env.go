package main

type environment struct {
	Port            string `env:"PUSH_NOTIFICATION_PORT,required"`
	Envrionment     string `env:"PUSH_NOTIFICATION_ENV,required"`
	ProjectID       string `env:"PROJECT_ID,required"`
	CredentialsPath string `env:"GOOGLE_APPLICATION_CREDENTIALS,required"`
	MinLogSeverity  string `env:"PUSH_NOTIFICATION_MIN_LOG_SEVERITY,required"`
	FcmServerKey    string `env:"FCM_SERVER_KEY,required"`
}
