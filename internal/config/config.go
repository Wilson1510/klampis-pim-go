type Config struct {
	App AppConfig
	Database DatabaseConfig
	JWT JWTConfig
}

type AppConfig struct {
	Name string
	Port string
	Debug bool
}

type DatabaseConfig struct {
	Host string
	Port string
	User string
	Password string
	Name string
}

type JWTConfig struct {
	Secret string
	AccessTokenExpiry time.Duration
}