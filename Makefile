# Variables
APP_NAME=todo-go
GO_FILES=$(shell find . -name "*.go" -type f)

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) main.go

run: ## Run the application locally
	@echo "Running $(APP_NAME)..."
	go run main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build files
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Database commands
migrate-up: ## Run database migrations
	@echo "Running migrations..."
	psql $(DATABASE_URL) -f database_schema.sql

migrate-down: ## Rollback database migrations (manual)
	@echo "Please manually drop tables or run specific rollback commands"

# Development commands
dev: ## Run in development mode with hot reload (requires air)
	@echo "Starting development server..."
	air

install-deps: ## Install Go dependencies
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest

# Production commands
deploy-build: ## Build for production
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/$(APP_NAME) main.go

# Linting and formatting
fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	golangci-lint run
