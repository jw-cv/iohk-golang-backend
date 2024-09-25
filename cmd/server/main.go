package main

import (
	"fmt"
	"iohk-golang-backend-preprod/internal/config"
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
}
