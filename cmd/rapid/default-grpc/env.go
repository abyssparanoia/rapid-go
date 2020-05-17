package defaultgrpc

type environment struct {
	Port            string `env:"DEFAULT_GRPC_PORT,required"`
	Envrionment     string `env:"DEFAULT_GRPC_ENV,required"`
	ProjectID       string `env:"PROJECT_ID,required"`
	CredentialsPath string `env:"GOOGLE_APPLICATION_CREDENTIALS,required"`
	DBHost          string `env:"DEFAULT_DB_HOST,required"`
	DBUser          string `env:"DEFAULT_DB_USER,required"`
	DBPassword      string `env:"DEFAULT_DB_PASSWORD,required"`
	DBDatabase      string `env:"DEFAULT_DB_DATABASE,required"`
}
