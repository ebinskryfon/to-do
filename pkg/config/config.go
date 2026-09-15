package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// DatabaseConfig holds connection and connection pooling settings for PostgreSQL.
type DatabaseConfig struct {
	URL             string
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}

// DSN returns the database connection string. If URL is explicitly set, it is returned;
// otherwise a connection URL is constructed from the individual parameters.
func (d DatabaseConfig) DSN() string {
	if d.URL != "" {
		return d.URL
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(d.User),
		url.QueryEscape(d.Password),
		d.Host,
		d.Port,
		d.Name,
		d.SSLMode,
	)
}

// Config holds all runtime configuration for the application.
type Config struct {
	AppEnv     string
	ServerPort string
	Database   DatabaseConfig

	// Flat field aliases for convenience and backwards compatibility
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// Load reads configuration from config/config.yaml (if present), a .env
// file (if present), and real environment variables, applying sane defaults.
func Load() (*Config, error) {
	// Load .env into the process environment if available
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./pkg/config")
	v.AddConfigPath(".")

	v.SetDefault("app_env", "development")
	v.SetDefault("server_port", "8080")

	v.SetDefault("database_url", "")
	v.SetDefault("db_host", "localhost")
	v.SetDefault("db_port", "5432")
	v.SetDefault("db_user", "postgres")
	v.SetDefault("db_password", "postgres")
	v.SetDefault("db_name", "todo_db")
	v.SetDefault("db_sslmode", "disable")

	// Connection pooling defaults (matching production-ready practices like Ollie Cake House)
	v.SetDefault("db_max_open_conns", 10)
	v.SetDefault("db_max_idle_conns", 2)
	v.SetDefault("db_conn_max_idle_time", "30s")
	v.SetDefault("db_conn_max_lifetime", "5m")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	idleTime, err := time.ParseDuration(v.GetString("db_conn_max_idle_time"))
	if err != nil {
		idleTime = 30 * time.Second
	}

	lifeTime, err := time.ParseDuration(v.GetString("db_conn_max_lifetime"))
	if err != nil {
		lifeTime = 5 * time.Minute
	}

	dbCfg := DatabaseConfig{
		URL:             v.GetString("database_url"),
		Host:            v.GetString("db_host"),
		Port:            v.GetString("db_port"),
		User:            v.GetString("db_user"),
		Password:        v.GetString("db_password"),
		Name:            v.GetString("db_name"),
		SSLMode:         v.GetString("db_sslmode"),
		MaxOpenConns:    v.GetInt("db_max_open_conns"),
		MaxIdleConns:    v.GetInt("db_max_idle_conns"),
		ConnMaxIdleTime: idleTime,
		ConnMaxLifetime: lifeTime,
	}

	cfg := &Config{
		AppEnv:     v.GetString("app_env"),
		ServerPort: v.GetString("server_port"),
		Database:   dbCfg,
		DBHost:     dbCfg.Host,
		DBPort:     dbCfg.Port,
		DBUser:     dbCfg.User,
		DBPassword: dbCfg.Password,
		DBName:     dbCfg.Name,
		DBSSLMode:  dbCfg.SSLMode,
	}

	return cfg, nil
}
