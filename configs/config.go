package configs

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Logger   LoggerConfig
	JWT      JWTConfig
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type AppConfig struct {
	Name string
	Host string
	Port int
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

type LoggerConfig struct {
	Level string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables: %v", err)
	}

	// Set defaults
	viper.SetDefault("APP_NAME", "my-gift")
	viper.SetDefault("APP_HOST", "0.0.0.0")
	viper.SetDefault("APP_PORT", 8080)
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_TIMEZONE", "UTC")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("JWT_SECRET", "change-me-in-production")
	viper.SetDefault("JWT_EXPIRY", "24h")

	jwtExpiry, err := time.ParseDuration(viper.GetString("JWT_EXPIRY"))
	if err != nil || jwtExpiry == 0 {
		jwtExpiry = 24 * time.Hour
	}

	cfg := &Config{
		App: AppConfig{
			Name: viper.GetString("APP_NAME"),
			Host: viper.GetString("APP_HOST"),
			Port: viper.GetInt("APP_PORT"),
			Env:  viper.GetString("APP_ENV"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
			TimeZone: viper.GetString("DB_TIMEZONE"),
		},
		Logger: LoggerConfig{
			Level: viper.GetString("LOG_LEVEL"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
			Expiry: jwtExpiry,
		},
	}

	return cfg, nil
}
