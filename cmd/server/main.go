package main

import (
	"context"
	"log"
	"net/http"

	"iohk-golang-backend-preprod/ent"
	"iohk-golang-backend-preprod/graph"
	"iohk-golang-backend-preprod/internal/config"
	"iohk-golang-backend-preprod/internal/domain/repository"
	"iohk-golang-backend-preprod/internal/domain/service"
	"iohk-golang-backend-preprod/internal/infra/db"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := loadConfig()
	pool := setupDatabasePool(cfg)
	defer db.CloseDBPool(pool)

	client := setupEntgoConnection(pool)
	defer client.Close()

	// Initialize repository and service
	customerRepo := repository.NewCustomerRepository(client)
	customerService := service.NewCustomerService(customerRepo)

	// Create resolver with the initialized service
	resolver := graph.NewResolver(customerService)

	// Set up GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Set up HTTP server
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://%s:%s/ for GraphQL playground", cfg.AppHost, cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, nil))
}

func loadConfig() *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return cfg
}

func setupDatabasePool(cfg *config.Config) *pgxpool.Pool {
	ctx := context.Background()
	pool, err := db.NewDBPool(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to set up database pool: %v", err)
	}
	return pool
}

func setupEntgoConnection(pool *pgxpool.Pool) *ent.Client {
	db := stdlib.OpenDBFromPool(pool)
	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))
	return client
}
