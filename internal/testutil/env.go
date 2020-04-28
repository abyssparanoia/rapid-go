package testutil

type environment struct {
	DefaultDBHost     string `env:"DEFAULT_DB_HOST,required"`
	DefaultDBUser     string `env:"DEFAULT_DB_USER,required"`
	DefaultDBDatabase string `env:"DEFAULT_DB_DATABASE,required"`
	DefaultDBPassword string `env:"DEFAULT_DB_PASSWORD,required"`
	DBLogEnabled      bool   `env:"DB_LOG_ENABLED" envDefault:"false"`
}
