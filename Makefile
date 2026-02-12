.PHONY: help build run test clean docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build -o digital-wallet .

run: ## Run the application
	go run main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -f digital-wallet
	go clean

deps: ## Download dependencies
	go mod download
	go mod tidy

docker-up: ## Start services with Docker Compose
	docker compose up -d

docker-down: ## Stop services
	docker compose down

docker-logs: ## Show Docker logs
	docker compose logs -f app

docker-rebuild: ## Rebuild and restart Docker services
	docker compose down
	docker compose build --no-cache
	docker compose up -d

up-postgre:
	@docker compose up -d postgres	

down-postgre:
	@docker compose down postgres

setup-db: ## Setup database (requires PostgreSQL running)
	psql -U postgres -d digital_wallet -f schema.sql

migrate: ## Run database migrations
	@echo "Running migrations..."
	psql -U postgres -d digital_wallet -f schema.sql

.DEFAULT_GOAL := help
