package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Arrange
	testCases := []struct {
		name           string
		envVars        map[string]string
		expectedConfig *Config
		expectedError  bool
	}{
		{
			name: "Valid configuration",
			envVars: map[string]string{
				"POSTGRES_USER":          "testuser",
				"POSTGRES_PASSWORD":      "testpass",
				"POSTGRES_DB":            "testdb",
				"POSTGRES_HOST":          "localhost",
				"POSTGRES_PORT":          "5432",
				"POSTGRES_SSLMODE":       "disable",
				"DB_MAX_CONNS":           "25",
				"DB_MIN_CONNS":           "5",
				"DB_MAX_CONN_LIFETIME":   "5h",
				"DB_MAX_CONN_IDLE_TIME":  "15m",
				"DB_HEALTH_CHECK_PERIOD": "1m",
				"APP_HOST":               "localhost",
				"APP_PORT":               "8080",
			},
			expectedConfig: &Config{
				PostgresUser:        "testuser",
				PostgresPassword:    "testpass",
				PostgresDB:          "testdb",
				PostgresHost:        "localhost",
				PostgresPort:        "5432",
				PostgresSSLMode:     "disable",
				DBMaxConns:          25,
				DBMinConns:          5,
				DBMaxConnLifetime:   5 * time.Hour,
				DBMaxConnIdleTime:   15 * time.Minute,
				DBHealthCheckPeriod: time.Minute,
				AppHost:             "localhost",
				AppPort:             "8080",
			},
			expectedError: false,
		},
		{
			name: "Partially set configuration",
			envVars: map[string]string{
				"POSTGRES_USER": "testuser",
			},
			expectedConfig: nil,
			expectedError:  true,
		},
		{
			name: "Empty configuration",
			envVars: map[string]string{
				"POSTGRES_USER":          "",
				"POSTGRES_PASSWORD":      "",
				"POSTGRES_DB":            "",
				"POSTGRES_HOST":          "",
				"POSTGRES_PORT":          "",
				"POSTGRES_SSLMODE":       "",
				"DB_MAX_CONNS":           "",
				"DB_MIN_CONNS":           "",
				"DB_MAX_CONN_LIFETIME":   "",
				"DB_MAX_CONN_IDLE_TIME":  "",
				"DB_HEALTH_CHECK_PERIOD": "",
				"APP_HOST":               "",
				"APP_PORT":               "",
			},
			expectedConfig: nil,
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			os.Clearenv()

			for key, value := range tc.envVars {
				os.Setenv(key, value)
			}

			// Act
			config, err := LoadConfig()

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
				assert.Equal(t, tc.expectedConfig, config)
			}
		})
	}
}

func TestLoadConfigWithoutEnvFile(t *testing.T) {
	// Arrange
	os.Clearenv()
	os.Setenv("POSTGRES_USER", "testuser")
	os.Setenv("POSTGRES_PASSWORD", "testpass")
	os.Setenv("POSTGRES_DB", "testdb")
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_SSLMODE", "disable")
	os.Setenv("DB_MAX_CONNS", "25")
	os.Setenv("DB_MIN_CONNS", "5")
	os.Setenv("DB_MAX_CONN_LIFETIME", "5h")
	os.Setenv("DB_MAX_CONN_IDLE_TIME", "15m")
	os.Setenv("DB_HEALTH_CHECK_PERIOD", "1m")
	os.Setenv("APP_HOST", "localhost")
	os.Setenv("APP_PORT", "8080")

	// Act
	config, err := LoadConfig()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "testuser", config.PostgresUser)
	assert.Equal(t, "testpass", config.PostgresPassword)
	assert.Equal(t, "testdb", config.PostgresDB)
	assert.Equal(t, "localhost", config.PostgresHost)
	assert.Equal(t, "5432", config.PostgresPort)
	assert.Equal(t, "disable", config.PostgresSSLMode)
	assert.Equal(t, 25, config.DBMaxConns)
	assert.Equal(t, 5, config.DBMinConns)
	assert.Equal(t, 5*time.Hour, config.DBMaxConnLifetime)
	assert.Equal(t, 15*time.Minute, config.DBMaxConnIdleTime)
	assert.Equal(t, time.Minute, config.DBHealthCheckPeriod)
	assert.Equal(t, "localhost", config.AppHost)
	assert.Equal(t, "8080", config.AppPort)
}

func TestValidateConfig(t *testing.T) {
	// Arrange
	testCases := []struct {
		name          string
		config        *Config
		expectedError string
	}{
		{
			name: "Valid configuration",
			config: &Config{
				PostgresUser:        "user",
				PostgresPassword:    "pass",
				PostgresDB:          "db",
				PostgresHost:        "host",
				PostgresPort:        "5432",
				PostgresSSLMode:     "disable",
				DBMaxConns:          25,
				DBMinConns:          5,
				DBMaxConnLifetime:   5 * time.Hour,
				DBMaxConnIdleTime:   15 * time.Minute,
				DBHealthCheckPeriod: time.Minute,
				AppHost:             "localhost",
				AppPort:             "8080",
			},
			expectedError: "",
		},
		{
			name: "Missing PostgresUser",
			config: &Config{
				PostgresPassword:    "pass",
				PostgresDB:          "db",
				PostgresHost:        "host",
				PostgresPort:        "5432",
				PostgresSSLMode:     "disable",
				DBMaxConns:          25,
				DBMinConns:          5,
				DBMaxConnLifetime:   5 * time.Hour,
				DBMaxConnIdleTime:   15 * time.Minute,
				DBHealthCheckPeriod: time.Minute,
				AppHost:             "localhost",
				AppPort:             "8080",
			},
			expectedError: "POSTGRES_USER is not set",
		},
		{
			name: "Missing DBMaxConns",
			config: &Config{
				PostgresUser:        "user",
				PostgresPassword:    "pass",
				PostgresDB:          "db",
				PostgresHost:        "host",
				PostgresPort:        "5432",
				PostgresSSLMode:     "disable",
				DBMinConns:          5,
				DBMaxConnLifetime:   5 * time.Hour,
				DBMaxConnIdleTime:   15 * time.Minute,
				DBHealthCheckPeriod: time.Minute,
				AppHost:             "localhost",
				AppPort:             "8080",
			},
			expectedError: "DB_MAX_CONNS must be greater than 0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			err := validateConfig(tc.config)

			// Assert
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
