.DEFAULT_GOAL := help
.PHONY: help run dev build generate swagger swag migrate migrate-up migrate-down bootstrap test tidy

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the API server locally
	go run ./cmd/api

dev: ## Run the API server with hot reload (requires air: go install github.com/air-verse/air@latest)
	air -c .air.toml

build: ## Build the API server binary into bin/api
	go build -o bin/api ./cmd/api

generate: ## Generate Swagger documentation from code annotations
	swag init -g cmd/api/main.go -o docs

swagger: generate ## Alias for generate
swag: generate ## Alias for generate

migrate: ## Ensure DB exists and run all pending migrations
	go run ./cmd/migrate up

migrate-up: ## Run all pending migrations
	go run ./cmd/migrate up

migrate-down: ## Rollback the latest migration
	go run ./cmd/migrate down

bootstrap: ## Seed initial development data (idempotent)
	go run ./cmd/bootstrap

test: ## Run all unit and package tests
	go test ./...

tidy: ## Clean up and download Go module dependencies
	go mod tidy
