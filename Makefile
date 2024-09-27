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

.PHONY: all build clean run test coverage lint vet fmt docker-build docker-up docker-down docker-logs help

all: build

build:
	@echo "Building..."
	@go build -o $(GOBIN)/$(BINARY_NAME) $(MAIN_PACKAGE)

clean:
	@echo "Cleaning..."
	@go clean
	@rm -f $(GOBIN)/$(BINARY_NAME)

run: build
	@echo "Running..."
	@$(GOBIN)/$(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test -v ./...

coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

lint:
	@echo "Linting..."
	@golangci-lint run

vet:
	@echo "Vetting..."
	@go vet ./...

fmt:
	@echo "Formatting..."
	@gofmt -s -w $(GO_FILES)

docker-build:
	@echo "Building Docker images..."
	@$(DOCKER_COMPOSE_CMD) build

docker-up:
	@echo "Starting Docker containers..."
	@$(DOCKER_COMPOSE_CMD) up -d

docker-down:
	@echo "Stopping Docker containers..."
	@$(DOCKER_COMPOSE_CMD) down

docker-logs:
	@echo "Showing Docker logs..."
	@$(DOCKER_COMPOSE_CMD) logs -f

generate:
	@echo "Generating GraphQL code..."
	@go run github.com/99designs/gqlgen generate

help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make clean          - Clean build files"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run tests"
	@echo "  make coverage       - Run tests with coverage"
	@echo "  make lint           - Run linter"
	@echo "  make vet            - Run go vet"
	@echo "  make fmt            - Format code"
	@echo "  make docker-build   - Build Docker images"
	@echo "  make docker-up      - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make docker-logs    - Show Docker logs"
	@echo "  make generate       - Generate GraphQL code"
	@echo "  make help           - Show this help message"
