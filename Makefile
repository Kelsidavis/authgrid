.PHONY: help build run dev stop clean test

# Detect docker compose command
DOCKER_COMPOSE := $(shell command -v docker-compose 2> /dev/null)
ifndef DOCKER_COMPOSE
	DOCKER_COMPOSE := docker compose
endif

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Docker containers
	$(DOCKER_COMPOSE) build

run: ## Start all services
	$(DOCKER_COMPOSE) up -d
	@echo ""
	@echo "✅ Authgrid is running!"
	@echo "   API:  http://localhost:8080"
	@echo "   Demo: http://localhost:3000"
	@echo ""
	@echo "View logs: make logs"
	@echo "Stop: make stop"

dev: ## Start services with logs
	$(DOCKER_COMPOSE) up

stop: ## Stop all services
	$(DOCKER_COMPOSE) down

clean: ## Stop services and remove volumes
	$(DOCKER_COMPOSE) down -v

logs: ## View logs
	$(DOCKER_COMPOSE) logs -f

logs-api: ## View API logs
	$(DOCKER_COMPOSE) logs -f api

logs-db: ## View database logs
	$(DOCKER_COMPOSE) logs -f postgres

db-shell: ## Open PostgreSQL shell
	$(DOCKER_COMPOSE) exec postgres psql -U authgrid -d authgrid

test: ## Run tests
	@echo "Running API tests..."
	@cd src/api && go test -v ./...
	@echo ""
	@echo "✅ All tests passed!"

install: ## Install Go dependencies
	cd src/api && go mod download

build-cli: ## Build CLI tool (using Docker)
	@echo "Building authgrid CLI..."
	@docker run --rm \
		-v "$(shell pwd)/src/cli":/src \
		-v "$(shell pwd)":/output \
		-w /src \
		golang:1.21-alpine \
		sh -c "go mod download && go build -o /output/authgrid ."
	@echo "✅ CLI built successfully: ./authgrid"
	@echo "Try: ./authgrid help"

install-cli: build-cli ## Install CLI to system
	@echo "Installing CLI to /usr/local/bin..."
	@sudo cp authgrid /usr/local/bin/
	@echo "✅ CLI installed! Try: authgrid help"

# Development helpers
.PHONY: db-reset db-migrate

db-reset: ## Reset database (WARNING: destroys all data)
	$(DOCKER_COMPOSE) down -v
	$(DOCKER_COMPOSE) up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Database reset complete"

db-migrate: ## Run database migrations manually
	$(DOCKER_COMPOSE) exec postgres psql -U authgrid -d authgrid -f /docker-entrypoint-initdb.d/001_init.sql
