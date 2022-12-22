package environment

type DatabaseEnvironment struct {
	DBHost     string `env:"DB_HOST,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBDatabase string `env:"DB_DATABASE,required"`
}

type Environment struct {
	Port        string `env:"PORT,required"`
	Environment string `env:"ENV,required"`
	DatabaseEnvironment

	GCPProjectID         string `env:"GCP_PROJECT_ID,required"`
	FirebaseClientAPIKey string `env:"FIREBASE_CLIENT_API_KEY"`
	BucketName           string `env:"BUCKET_NAME,required"`
}
