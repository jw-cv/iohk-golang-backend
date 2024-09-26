package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"iohk-golang-backend-preprod/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	// TODO: put sslmode in env variables
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database connection string: %w", err)
	}

	// TODO: put this in env variables
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 5 * time.Hour
	poolConfig.MaxConnIdleTime = 15 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to the database")
	return pool, nil
}

func CloseDBPool(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
		log.Println("Database connection pool closed")
	}
}

func SetupDBPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := NewDBPool(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}
	return pool, nil
}

// TODO: Use the pool in your handlers or resolvers:

// func someHandler(pool *pgxpool.Pool) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 			ctx := r.Context()
// 			// Use pool.Acquire() to get a connection from the pool
// 			conn, err := pool.Acquire(ctx)
// 			if err != nil {
// 					http.Error(w, "Failed to acquire database connection", http.StatusInternalServerError)
// 					return
// 			}
// 			defer conn.Release()

// 			// Use conn.Query(), conn.QueryRow(), or conn.Exec() for database operations
// 			// ...
// 	}
// }
