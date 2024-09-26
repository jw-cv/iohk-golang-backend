package config

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresUser        string
	PostgresPassword    string
	PostgresDB          string
	PostgresHost        string
	PostgresPort        string
	PostgresSSLMode     string
	DBMaxConns          int
	DBMinConns          int
	DBMaxConnLifetime   time.Duration
	DBMaxConnIdleTime   time.Duration
	DBHealthCheckPeriod time.Duration
	AppPort             string
}

func LoadConfig() (*Config, error) {
	envFile := ".env.local"

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Error loading %s file: %v", envFile, err)
	}

	viper.AutomaticEnv()

	config := &Config{
		PostgresUser:        viper.GetString("POSTGRES_USER"),
		PostgresPassword:    viper.GetString("POSTGRES_PASSWORD"),
		PostgresDB:          viper.GetString("POSTGRES_DB"),
		PostgresHost:        viper.GetString("POSTGRES_HOST"),
		PostgresPort:        viper.GetString("POSTGRES_PORT"),
		PostgresSSLMode:     viper.GetString("POSTGRES_SSLMODE"),
		DBMaxConns:          viper.GetInt("DB_MAX_CONNS"),
		DBMinConns:          viper.GetInt("DB_MIN_CONNS"),
		DBMaxConnLifetime:   viper.GetDuration("DB_MAX_CONN_LIFETIME"),
		DBMaxConnIdleTime:   viper.GetDuration("DB_MAX_CONN_IDLE_TIME"),
		DBHealthCheckPeriod: viper.GetDuration("DB_HEALTH_CHECK_PERIOD"),
		AppPort:             viper.GetString("APP_PORT"),
	}

	if err := validateConfig(config); err != nil {
		return nil, err // Return nil and the error if validation fails
	}

	return config, nil
}

func validateConfig(c *Config) error {
	if c.PostgresUser == "" {
		return fmt.Errorf("POSTGRES_USER is not set")
	}
	if c.PostgresPassword == "" {
		return fmt.Errorf("POSTGRES_PASSWORD is not set")
	}
	if c.PostgresDB == "" {
		return fmt.Errorf("POSTGRES_DB is not set")
	}
	if c.PostgresHost == "" {
		return fmt.Errorf("POSTGRES_HOST is not set")
	}
	if c.PostgresPort == "" {
		return fmt.Errorf("POSTGRES_PORT is not set")
	}
	if c.PostgresSSLMode == "" {
		return fmt.Errorf("POSTGRES_SSLMODE is not set")
	}
	if c.DBMaxConns <= 0 {
		return fmt.Errorf("DB_MAX_CONNS must be greater than 0")
	}
	if c.DBMinConns <= 0 {
		return fmt.Errorf("DB_MIN_CONNS must be greater than 0")
	}
	if c.DBMaxConnLifetime <= 0 {
		return fmt.Errorf("DB_MAX_CONN_LIFETIME must be greater than 0")
	}
	if c.DBMaxConnIdleTime <= 0 {
		return fmt.Errorf("DB_MAX_CONN_IDLE_TIME must be greater than 0")
	}
	if c.DBHealthCheckPeriod <= 0 {
		return fmt.Errorf("DB_HEALTH_CHECK_PERIOD must be greater than 0")
	}
	if c.AppPort == "" {
		return fmt.Errorf("APP_PORT is not set")
	}
	return nil
}
