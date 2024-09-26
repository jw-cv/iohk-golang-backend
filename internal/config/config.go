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
	fields := map[string]string{
		"POSTGRES_USER":          c.PostgresUser,
		"POSTGRES_PASSWORD":      c.PostgresPassword,
		"POSTGRES_DB":            c.PostgresDB,
		"POSTGRES_HOST":          c.PostgresHost,
		"POSTGRES_PORT":          c.PostgresPort,
		"POSTGRES_SSLMODE":       c.PostgresSSLMode,
		"DB_MAX_CONNS":           fmt.Sprintf("%d", c.DBMaxConns),
		"DB_MIN_CONNS":           fmt.Sprintf("%d", c.DBMinConns),
		"DB_MAX_CONN_LIFETIME":   fmt.Sprintf("%s", c.DBMaxConnLifetime),
		"DB_MAX_CONN_IDLE_TIME":  fmt.Sprintf("%s", c.DBMaxConnIdleTime),
		"DB_HEALTH_CHECK_PERIOD": fmt.Sprintf("%s", c.DBHealthCheckPeriod),
		"APP_PORT":               c.AppPort,
	}

	for key, value := range fields {
		if value == "" {
			return fmt.Errorf("%s is not set", key)
		}
	}

	return nil
}
