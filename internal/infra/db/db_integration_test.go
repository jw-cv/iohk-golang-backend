//go:build integration
// +build integration

package db_test

import (
	"context"
	"os"
	"testing"
	"time"

	"iohk-golang-backend-preprod/internal/config"
	"iohk-golang-backend-preprod/internal/infra/db"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Load .env.test file
	err := godotenv.Load("../../../.env.test")
	require.NoError(t, err, "Error loading .env.test file")

	// Override POSTGRES_HOST to use localhost
	os.Setenv("POSTGRES_HOST", "localhost")

	cfg, err := config.LoadConfig()
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := db.NewDBPool(ctx, cfg)
	require.NoError(t, err)
	defer pool.Close()

	err = pool.Ping(ctx)
	assert.NoError(t, err)

	// Add more specific tests here
	t.Run("Query customers_test", func(t *testing.T) {
		var count int
		err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM customers_test").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 2, count, "Expected 2 customers in the test database")
	})
}
