package environment

type Environment struct {
	Port        string `env:"PORT,required"`
	Environment string `env:"ENV,required"`
	DatabaseEnvironment
	GCPEnviroment
	AWSEnvironment
}

type DatabaseEnvironment struct {
	DBHost     string `env:"DB_HOST,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBDatabase string `env:"DB_DATABASE,required"`
}

type GCPEnviroment struct {
	GCPProjectID         string `env:"GCP_PROJECT_ID,required"`
	FirebaseClientAPIKey string `env:"FIREBASE_CLIENT_API_KEY"`
	GCPBucketName        string `env:"GCP_BUCKET_NAME,required"`
}

type AWSEnvironment struct {
	AWSBucketName   string `env:"AWS_BUCKET_NAME,required"`
	AWSRegion       string `env:"AWS_REGION,required"`
	AWSEmulatorHost string `env:"AWS_EMULATOR_HOST"`
}
