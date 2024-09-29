# Variables
BINARY_NAME=iohk-golang-backend-preprod
DOCKER_COMPOSE_FILE=docker-compose.yml
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")
ENV_FILE=.env.local

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

# Main package path
MAIN_PACKAGE=./cmd/server

# Determine Docker Compose command
DOCKER_COMPOSE := $(shell \
    if command -v docker-compose >/dev/null 2>&1; then \
        echo "docker-compose"; \
    elif docker compose version >/dev/null 2>&1; then \
        echo "docker compose"; \
    else \
        echo ""; \
    fi)

# Check if Docker Compose is available
ifeq ($(DOCKER_COMPOSE),)
    $(error "Neither 'docker-compose' nor 'docker compose' found. Please install Docker Compose.")
else
    DOCKER_COMPOSE_CMD := $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE)
endif

# Ensure GOPATH is set before running build
GOPATH ?= $(HOME)/go

.PHONY: all build clean run test coverage test-integration-up test-integration-run test-integration-down test-integration lint vet fmt docker-build docker-up docker-down docker-logs help

all: build

# Docker-related commands
docker-build:
	@echo "Building Docker images..."
	@$(DOCKER_COMPOSE_CMD) build

docker-up:
	@echo "Starting Docker containers..."
	@$(DOCKER_COMPOSE_CMD) up -d
	@echo "Showing Docker logs..."
	@$(DOCKER_COMPOSE_CMD) logs -f

docker-down:
	@echo "Stopping Docker containers..."
	@$(DOCKER_COMPOSE_CMD) down -v

docker-logs:
	@echo "Showing Docker logs..."
	@$(DOCKER_COMPOSE_CMD) logs -f

# Local development commands
clean:
	@echo "Cleaning..."
	@go clean
	@rm -f $(GOBIN)/$(BINARY_NAME)

build:
	@echo "Ensuring dependencies are downloaded..."
	@go mod download
	@echo "Building the application locally without docker (for development purposes only)..."
	@go build -o $(GOBIN)/$(BINARY_NAME) $(MAIN_PACKAGE)
run: build
	@echo "Running the application locally without docker (for development purposes)..."
	@$(GOBIN)/$(BINARY_NAME)

# Test related commands
test:
	@echo "Ensuring dependencies are downloaded..."
	@go mod download
	@echo "Running tests..."
	@go test -v ./...

coverage:
	@echo "Ensuring dependencies are downloaded..."
	@go mod download
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Integration tests commands
test-integration-up:
	@echo "Starting test database..."
	@$(DOCKER_COMPOSE_CMD) -f docker-compose.test.yml up -d
	@echo "Waiting for test database to be ready..."
	@sleep 5

test-integration-run:
	@echo "Running integration tests..."
	@go test -v -tags integration ./... -run TestDatabaseConnection

test-integration-down:
	@echo "Stopping test database..."
	@$(DOCKER_COMPOSE_CMD) -f docker-compose.test.yml down -v

test-integration: test-integration-up test-integration-run test-integration-down

lint:
	@echo "Ensuring dependencies are downloaded..."
	@go mod download
	@echo "Linting..."
	@golangci-lint run

vet:
	@echo "Ensuring dependencies are downloaded..."
	@go mod download
	@echo "Vetting..."
	@go vet ./...

fmt:
	@echo "Ensuring dependencies are downloaded..."
	@go mod download
	@echo "Formatting..."
	@gofmt -s -w $(GO_FILES)

generate:
	@echo "Ensuring dependencies are downloaded..."
	@go mod download
	@echo "Generating GraphQL code..."
	@go run github.com/99designs/gqlgen generate

help:
	@echo "Available commands:"
	@echo "  make docker-build         - Build Docker images"
	@echo "  make docker-up            - Start Docker containers"
	@echo "  make docker-down          - Stop Docker containers"
	@echo "  make docker-logs          - Show Docker logs"
	@echo "  make clean                - Clean build files"
	@echo "  make build                - Build the application locally without docker (for development purposes only)"
	@echo "  make run                  - Run the application locally without docker (for development purposes only)"
	@echo "  make test                 - Run unit tests"
	@echo "  make coverage             - Run tests with coverage"
	@echo "  make test-integration-up   - Start the test database for integration tests"
	@echo "  make test-integration-run  - Run integration tests (assumes test DB is already up)"
	@echo "  make test-integration-down - Stop the test database for integration tests"
	@echo "  make test-integration      - Run full integration test cycle (up, test, down)"
	@echo "  make lint                 - Run linter"
	@echo "  make vet                  - Run go vet"
	@echo "  make fmt                  - Format code"
	@echo "  make generate             - Generate GraphQL code"
	@echo "  make help                 - Show this help message"