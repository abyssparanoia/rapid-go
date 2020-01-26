package gluepsql

// Config ... psql config
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// NewConfig ... get new config
func NewConfig(host string, port int, user string, password string, database string) *Config {
	return &Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
	}
}
