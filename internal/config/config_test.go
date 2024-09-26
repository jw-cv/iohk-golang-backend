package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {

	testCases := []struct {
		name           string
		envVars        map[string]string
		expectedConfig *Config
		expectedError  bool
	}{
		{
			name: "Valid configuration",
			envVars: map[string]string{
				"POSTGRES_USER":     "testuser",
				"POSTGRES_PASSWORD": "testpass",
				"POSTGRES_DB":       "testdb",
				"POSTGRES_HOST":     "localhost",
				"POSTGRES_PORT":     "5432",
				"APP_PORT":          "8080",
			},
			expectedConfig: &Config{
				PostgresUser:     "testuser",
				PostgresPassword: "testpass",
				PostgresDB:       "testdb",
				PostgresHost:     "localhost",
				PostgresPort:     "5432",
				AppPort:          "8080",
			},
			expectedError: false,
		},
		{
			name: "Missing configuration",
			envVars: map[string]string{
				"POSTGRES_USER": "testuser",
				// Missing other required fields
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
	assert.Equal(t, "8080", config.AppPort)
}

func TestValidateConfig(t *testing.T) {
	// Arrange
	testCases := []struct {
		name          string
		config        *Config
		expectedError bool
	}{
		{
			name: "Valid configuration",
			config: &Config{
				PostgresUser:     "user",
				PostgresPassword: "pass",
				PostgresDB:       "db",
				PostgresHost:     "host",
				PostgresPort:     "5432",
				AppPort:          "8080",
			},
			expectedError: false,
		},
		{
			name: "Missing PostgresUser",
			config: &Config{
				PostgresPassword: "pass",
				PostgresDB:       "db",
				PostgresHost:     "host",
				PostgresPort:     "5432",
				AppPort:          "8080",
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			err := validateConfig(tc.config)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
