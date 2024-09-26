package main

import (
	"context"
	"fmt"
	"iohk-golang-backend-preprod/internal/config"
	"iohk-golang-backend-preprod/internal/infra/db"
	"log"
)

func main() {
	// Entry point of the application
	fmt.Println("Hello World")

	// Load the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Use the configuration values
	fmt.Printf("PostgresUser: %s\n", cfg.PostgresUser)
	fmt.Printf("PostgresPassword: %s\n", cfg.PostgresPassword)
	fmt.Printf("PostgresDB: %s\n", cfg.PostgresDB)
	fmt.Printf("PostgresHost: %s\n", cfg.PostgresHost)
	fmt.Printf("PostgresPort: %s\n", cfg.PostgresPort)
	fmt.Printf("AppPort: %s\n", cfg.AppPort)

	// Set up the database connection
	ctx := context.Background()
	pool, err := db.SetupDBPool(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to set up database pool: %v", err)
	}
	defer db.CloseDBPool(pool)

	log.Println("Database connection established")

	// Set up your web server or GraphQL API here
	// ...

	// Run your server
	// ...
}
