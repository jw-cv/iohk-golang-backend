package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"iohk-golang-backend-preprod/ent"
	"iohk-golang-backend-preprod/ent/customer"
	"iohk-golang-backend-preprod/ent/migrate"
	"iohk-golang-backend-preprod/graph"
	"iohk-golang-backend-preprod/internal/config"
	"iohk-golang-backend-preprod/internal/infra/db"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jackc/pgx/v5/stdlib"
)

const defaultPort = "8080"

func main() {

	// Load the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up the database connection
	ctx := context.Background()
	pool, err := db.NewDBPool(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to set up database pool: %v", err)
	}
	defer db.CloseDBPool(pool)

	// Set up entgo connection
	db := stdlib.OpenDBFromPool(pool)
	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))

	CreateCustomer(context.Background(), client)
	QueryCustomer(context.Background(), client)
	AutoMigration(context.Background(), client)

	// Set up your web server or GraphQL API here
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	// Run your server
	// ...
}

func CreateCustomer(ctx context.Context, client *ent.Client) (*ent.Customer, error) {
	customer, err := client.Customer.
		Create().
		SetName("Ron Burgandy").
		SetSurname("Burgandy").
		SetNumber(333).
		SetGender("Male").
		SetCountry("USA").
		SetDependants(69).
		SetBirthDate(time.Date(1990, time.February, 2, 0, 0, 0, 0, time.UTC)).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating customer: %w", err)
	}
	log.Println("customer was created: ", customer)
	return customer, nil
}

func QueryCustomer(ctx context.Context, client *ent.Client) (*ent.Customer, error) {
	customer, err := client.Customer.
		Query().
		Where(customer.Name("Ron Burgandy")).
		// `Only` fails if no customer found,
		// or more than 1 customer returned.
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying customer: %w", err)
	}
	log.Println("customer returned: ", customer)
	return customer, nil
}

func AutoMigration(ctx context.Context, client *ent.Client) {
	log.Println("creating new schema resources (auto migration)")
	err := client.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
		migrate.WithGlobalUniqueID(true),
	)
	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	log.Println("schema resources created successfully (auto migration)")
}
