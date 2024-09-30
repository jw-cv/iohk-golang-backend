//go:build integration && testcoverage
// +build integration,testcoverage

package db_test

import (
	"context"
	"testing"
	"time"

	"iohk-golang-backend/internal/config"
	"iohk-golang-backend/internal/domain/model"
	"iohk-golang-backend/internal/infra/db"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestCustomer struct {
	ID         int
	Name       string
	Surname    string
	Number     int
	Gender     string
	Country    string
	Dependants int
	BirthDate  time.Time
}

func TestDatabaseIntegration(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)
	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Create a test configuration
	cfg := &config.Config{
		PostgresUser:        "testuser",
		PostgresPassword:    "testpass",
		PostgresDB:          "testdb",
		PostgresHost:        host,
		PostgresPort:        port.Port(),
		PostgresSSLMode:     "disable",
		DBMaxConns:          5,
		DBMinConns:          1,
		DBMaxConnLifetime:   time.Hour,
		DBMaxConnIdleTime:   time.Minute * 30,
		DBHealthCheckPeriod: time.Minute,
	}

	// Create the database pool
	pool, err := db.NewDBPool(ctx, cfg)
	require.NoError(t, err)
	defer pool.Close()

	// Run the tests
	t.Run("CreateSchema", testCreateSchema(ctx, pool))
	t.Run("InsertCustomer", testInsertCustomer(ctx, pool))
	t.Run("GetCustomer", testGetCustomer(ctx, pool))
}

func testCreateSchema(ctx context.Context, pool *pgxpool.Pool) func(*testing.T) {
	return func(t *testing.T) {
		_, err := pool.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS customers (
				id SERIAL PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				surname VARCHAR(255) NOT NULL,
				number INTEGER NOT NULL,
				gender VARCHAR(10) CHECK (gender IN ('Male', 'Female')),
				country VARCHAR(255),
				dependants INTEGER CHECK (dependants >= 0),
				birth_date DATE CHECK (birth_date <= CURRENT_DATE)
			)
		`)
		assert.NoError(t, err)
	}
}

func testInsertCustomer(ctx context.Context, pool *pgxpool.Pool) func(*testing.T) {
	return func(t *testing.T) {
		customer := TestCustomer{
			Name:       "John",
			Surname:    "Doe",
			Number:     123,
			Gender:     "Male",
			Country:    "USA",
			Dependants: 2,
			BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			ID:         1, // Assuming ID is 1 for simplicity
		}

		_, err := pool.Exec(ctx, `
			INSERT INTO customers (name, surname, number, gender, country, dependants, birth_date)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, customer.Name, customer.Surname, customer.Number, customer.Gender,
			customer.Country, customer.Dependants, customer.BirthDate)
		assert.NoError(t, err)
	}
}

func testGetCustomer(ctx context.Context, pool *pgxpool.Pool) func(*testing.T) {
	return func(t *testing.T) {
		var customer TestCustomer
		err := pool.QueryRow(ctx, `
			SELECT id, name, surname, number, gender, country, dependants, birth_date
			FROM customers
			WHERE name = $1 AND surname = $2
		`, "John", "Doe").Scan(
			&customer.ID, &customer.Name, &customer.Surname, &customer.Number,
			&customer.Gender, &customer.Country, &customer.Dependants, &customer.BirthDate)
		assert.NoError(t, err)

		assert.Equal(t, "John", customer.Name)
		assert.Equal(t, "Doe", customer.Surname)
		assert.Equal(t, 123, customer.Number)
		assert.Equal(t, "Male", customer.Gender)
		assert.Equal(t, "USA", customer.Country)
		assert.Equal(t, 2, customer.Dependants)
		assert.Equal(t, time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), customer.BirthDate)

		// Convert the gender string to the model.Gender type
		domainGender := model.Gender(customer.Gender)
		assert.Equal(t, model.Gender("Male"), domainGender)
	}
}
