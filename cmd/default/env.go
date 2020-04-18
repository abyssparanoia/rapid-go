package main

type environment struct {
	Port            string `env:"DEFAULT_PORT,required"`
	Envrionment     string `env:"DEFAULT_ENV,required"`
	ProjectID       string `env:"PROJECT_ID,required"`
	CredentialsPath string `env:"GOOGLE_APPLICATION_CREDENTIALS,required"`
	MinLogSeverity  string `env:"DEFAULT_MIN_LOG_SEVERITY,required"`
	DBHost          string `env:"DEFAULT_DB_HOST,required"`
	DBUser          string `env:"DEFAULT_DB_USER,required"`
	DBPassword      string `env:"DEFAULT_DB_PASSWORD,required"`
	DBDatabase      string `env:"DEFAULT_DB_DATABASE,required"`
}
