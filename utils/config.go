package utils

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Configuration struct {
	AppName     string
	Port        string
	Debug       bool
	Limit       int
	PathLogging string
	DB          DatabaseConfig
	JWT         JWTConfig
	Redis       RedisConfig
}

type DatabaseConfig struct {
	Name     string
	Username string
	Password string
	Host     string
	Port     string
	MaxConn  int32
}

// Backwards compatibility alias
type DatabaseCofig = DatabaseConfig

type JWTConfig struct {
	Secret        string
	ExpiryMinutes time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Enabled  bool
}

// ReadConfiguration reads configuration from .env and environment variables
func ReadConfiguration() (Configuration, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Set default values
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("APP_NAME", "AppInventory")
	viper.SetDefault("DEBUG", true)
	viper.SetDefault("LIMIT", 10)
	viper.SetDefault("PATH_LOGGING", "./logs/")
	viper.SetDefault("DATABASE_HOST", "localhost")
	viper.SetDefault("DATABASE_PORT", "5432")
	viper.SetDefault("DATABASE_MAX_CONN", 10)
	viper.SetDefault("JWT_SECRET", "super-secret-inventory-jwt-key-2026")
	viper.SetDefault("JWT_EXPIRY_MINUTES", 1440) // 24 hours
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("REDIS_ENABLED", false)

	_ = viper.ReadInConfig() // ignore error if .env doesn't exist, fallback to env vars

	viper.AutomaticEnv()

	pflag.String("port-app", "", "port for app golang")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	appPort := viper.GetString("PORT")
	if flagPort := viper.GetString("port-app"); flagPort != "" {
		appPort = flagPort
	}
	if appPort == "" {
		appPort = "8080"
	}

	limit := viper.GetInt("LIMIT")
	if limit <= 0 {
		limit = 10
	}

	jwtExpiry := time.Duration(viper.GetInt("JWT_EXPIRY_MINUTES")) * time.Minute
	if jwtExpiry <= 0 {
		jwtExpiry = 24 * time.Hour
	}

	return Configuration{
		AppName:     viper.GetString("APP_NAME"),
		Port:        appPort,
		Debug:       viper.GetBool("DEBUG"),
		Limit:       limit,
		PathLogging: viper.GetString("PATH_LOGGING"),
		DB: DatabaseConfig{
			Name:     viper.GetString("DATABASE_NAME"),
			Username: viper.GetString("DATABASE_USERNAME"),
			Password: viper.GetString("DATABASE_PASSWORD"),
			Host:     viper.GetString("DATABASE_HOST"),
			Port:     viper.GetString("DATABASE_PORT"),
			MaxConn:  viper.GetInt32("DATABASE_MAX_CONN"),
		},
		JWT: JWTConfig{
			Secret:        viper.GetString("JWT_SECRET"),
			ExpiryMinutes: jwtExpiry,
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
			Enabled:  viper.GetBool("REDIS_ENABLED"),
		},
	}, nil
}
