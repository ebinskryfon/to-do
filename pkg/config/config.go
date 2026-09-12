package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all runtime configuration for the application, loaded from
// config.yaml and/or environment variables (env vars take precedence).
type Config struct {
	AppEnv     string
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// Load reads configuration from config/config.yaml (if present), a .env
// file (if present), and real environment variables, applying sane
// local-dev defaults. Precedence, highest first: real env vars > .env >
// config.yaml > defaults.
func Load() (*Config, error) {
	// Load .env into the process environment so AutomaticEnv() below picks
	// it up. Ignored if the file doesn't exist — .env is optional, e.g. in
	// production where real env vars are set directly.
	_ = godotenv.Load()

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./pkg/config")

	v.SetDefault("app_env", "development")
	v.SetDefault("server_port", "8080")
	v.SetDefault("db_host", "localhost")
	v.SetDefault("db_port", "5432")
	v.SetDefault("db_user", "postgres")
	v.SetDefault("db_password", "postgres")
	v.SetDefault("db_name", "todo_db")
	v.SetDefault("db_sslmode", "disable")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// config.yaml is optional — env vars/.env alone are enough to run.
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	cfg := &Config{
		AppEnv:     v.GetString("app_env"),
		ServerPort: v.GetString("server_port"),
		DBHost:     v.GetString("db_host"),
		DBPort:     v.GetString("db_port"),
		DBUser:     v.GetString("db_user"),
		DBPassword: v.GetString("db_password"),
		DBName:     v.GetString("db_name"),
		DBSSLMode:  v.GetString("db_sslmode"),
	}

	return cfg, nil
}
