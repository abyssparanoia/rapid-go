package environment

type Environment struct {
	Port        string `env:"PORT,required"`
	Environment string `env:"ENV,required"`
	DatabaseEnvironment
	RedisEnvironment
	GCPEnvironment
	AWSEnvironment
	SpannerEnvironment
}

type DatabaseEnvironment struct {
	DBHost      string `env:"DB_HOST,required"`
	DBUser      string `env:"DB_USER,required"`
	DBPassword  string `env:"DB_PASSWORD,required"`
	DBDatabase  string `env:"DB_DATABASE,required"`
	DBLogEnable bool   `env:"DB_LOG_ENABLE"        envDefault:"false"`
}

type RedisEnvironment struct {
	RedisHost      string `env:"REDIS_HOST,required"`
	RedisPort      string `env:"REDIS_PORT,required"`
	RedisUsername  string `env:"REDIS_USERNAME"`
	RedisPassword  string `env:"REDIS_PASSWORD,required"`
	RedisTLSEnable bool   `env:"REDIS_TLS_ENABLE"        envDefault:"true"`
}

type GCPEnvironment struct {
	GCPProjectID         string `env:"GCP_PROJECT_ID,required"`
	FirebaseClientAPIKey string `env:"FIREBASE_CLIENT_API_KEY"`
	GCPBucketName        string `env:"GCP_BUCKET_NAME,required"`
}

type AWSEnvironment struct {
	AWSRegion                 string `env:"AWS_REGION,required"`
	AWSEmulatorHost           string `env:"AWS_EMULATOR_HOST"`
	AWSCognitoEmulatorHost    string `env:"AWS_COGNITO_EMULATOR_HOST"`
	AWSBucketName             string `env:"AWS_BUCKET_NAME,required"`
	AWSCognitoStaffUserPoolID string `env:"AWS_COGNITO_STAFF_USER_POOL_ID,required"`
	AWSCognitoStaffClientID   string `env:"AWS_COGNITO_STAFF_CLIENT_ID,required"`
}

type SpannerEnvironment struct {
	SpannerProjectID  string `env:"SPANNER_PROJECT_ID,required"`
	SpannerInstanceID string `env:"SPANNER_INSTANCE_ID,required"`
	SpannerDatabaseID string `env:"SPANNER_DATABASE_ID,required"`
}
