# Variables
MAKEFLAGS += -j4

# Run the application in development mode
dev:
	skaffold dev --default-repo registry.localhost:5000 --profile dev

# Run the application in production mode
start:
	skaffold run

# Stop
stop:
	skaffold delete --profile dev

# K8S Cluster Commands
cluster-start:
	k3d cluster create tower --config k3d.yaml

cluster-stop:
	k3d cluster delete -a

test-e2e:
	skaffold build --file-output=build.json --profile dev
	skaffold exec test-e2e --build-artifacts=build.json --profile dev

cluster-upgrade:
	helm upgrade dev ./helm/arches -f ./helm/dev-overrides.yaml

cluster-install:
	helm install dev ./helm/arches -f ./helm/dev-overrides.yaml



# .PHONY: help build run test clean migrate generate sqlc oapi dev lint fmt

# # Variables
# BINARY_NAME=archesai-api
# MAIN_PATH=cmd/api/main.go
# MIGRATION_PATH=internal/database/migrations
# DATABASE_URL ?= postgres://localhost/archesai?sslmode=disable

# # Colors for output
# GREEN=\033[0;32m
# YELLOW=\033[0;33m
# NC=\033[0m # No Color

# help: ## Show this help message
# 	@echo 'Usage: make [target]'
# 	@echo ''
# 	@echo 'Available targets:'
# 	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${GREEN}%-15s${NC} %s\n", $$1, $$2}'  $(MAKEFILE_LIST)

# build: ## Build the application
# 	@echo "${YELLOW}Building application...${NC}"
# 	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)
# 	@echo "${GREEN}Build complete!${NC}"

# run: ## Run the application
# 	@echo "${YELLOW}Running application...${NC}"
# 	go run $(MAIN_PATH)

# dev: ## Run the application with hot reload (requires air)
# 	@echo "${YELLOW}Running in development mode...${NC}"
# 	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
# 	air

# test: ## Run tests
# 	@echo "${YELLOW}Running tests...${NC}"
# 	go test -v -cover ./...
# 	@echo "${GREEN}Tests complete!${NC}"

# test-coverage: ## Run tests with coverage report
# 	@echo "${YELLOW}Running tests with coverage...${NC}"
# 	go test -v -coverprofile=coverage.out ./...
# 	go tool cover -html=coverage.out -o coverage.html
# 	@echo "${GREEN}Coverage report generated: coverage.html${NC}"

# clean: ## Clean build artifacts
# 	@echo "${YELLOW}Cleaning...${NC}"
# 	rm -rf bin/
# 	rm -f coverage.out coverage.html
# 	@echo "${GREEN}Clean complete!${NC}"

# deps: ## Download dependencies
# 	@echo "${YELLOW}Downloading dependencies...${NC}"
# 	go mod download
# 	go mod tidy
# 	@echo "${GREEN}Dependencies updated!${NC}"

# generate: sqlc oapi ## Generate all code (sqlc + OpenAPI)
# 	@echo "${GREEN}All code generation complete!${NC}"

# sqlc: ## Generate database code with sqlc
# 	@echo "${YELLOW}Generating sqlc code...${NC}"
# 	@which sqlc > /dev/null || (echo "Installing sqlc..." && go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest)
# 	sqlc generate
# 	@echo "${GREEN}sqlc generation complete!${NC}"

# oapi: ## Generate server code from OpenAPI spec
# 	@echo "${YELLOW}Generating OpenAPI server code...${NC}"
# 	@which oapi-codegen > /dev/null || (echo "Installing oapi-codegen..." && go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest)
# 	mkdir -p internal/api/generated
# 	oapi-codegen -package generated -generate types,server,spec api/openapi.yaml > internal/api/generated/server.gen.go
# 	@echo "${GREEN}OpenAPI generation complete!${NC}"

# migrate-up: ## Run database migrations up
# 	@echo "${YELLOW}Running migrations up...${NC}"
# 	@which migrate > /dev/null || (echo "Please install golang-migrate: https://github.com/golang-migrate/migrate" && exit 1)
# 	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up
# 	@echo "${GREEN}Migrations complete!${NC}"

# migrate-down: ## Run database migrations down
# 	@echo "${YELLOW}Running migrations down...${NC}"
# 	@which migrate > /dev/null || (echo "Please install golang-migrate: https://github.com/golang-migrate/migrate" && exit 1)
# 	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down
# 	@echo "${GREEN}Migrations rolled back!${NC}"

# migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
# 	@echo "${YELLOW}Creating migration: $(name)...${NC}"
# 	@which migrate > /dev/null || (echo "Please install golang-migrate: https://github.com/golang-migrate/migrate" && exit 1)
# 	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)
# 	@echo "${GREEN}Migration created!${NC}"

# lint: ## Run linter
# 	@echo "${YELLOW}Running linter...${NC}"
# 	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
# 	golangci-lint run ./...
# 	@echo "${GREEN}Linting complete!${NC}"

# fmt: ## Format code
# 	@echo "${YELLOW}Formatting code...${NC}"
# 	go fmt ./...
# 	goimports -w .
# 	@echo "${GREEN}Formatting complete!${NC}"

# docker-build: ## Build Docker image
# 	@echo "${YELLOW}Building Docker image...${NC}"
# 	docker build -t $(BINARY_NAME):latest .
# 	@echo "${GREEN}Docker build complete!${NC}"

# docker-run: ## Run application in Docker
# 	@echo "${YELLOW}Running Docker container...${NC}"
# 	docker run -p 8080:8080 --env-file .env $(BINARY_NAME):latest

# docker-compose-up: ## Start services with docker-compose
# 	@echo "${YELLOW}Starting services...${NC}"
# 	docker-compose up -d
# 	@echo "${GREEN}Services started!${NC}"

# docker-compose-down: ## Stop services with docker-compose
# 	@echo "${YELLOW}Stopping services...${NC}"
# 	docker-compose down
# 	@echo "${GREEN}Services stopped!${NC}"

# install-tools: ## Install development tools
# 	@echo "${YELLOW}Installing development tools...${NC}"
# 	go install github.com/air-verse/air@latest
# 	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
# 	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
# 	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
# 	go install golang.org/x/tools/cmd/goimports@latest
# 	@echo "${GREEN}Tools installed!${NC}"







# Simple Makefile for a Go project

# Build the application
# all: build test

# build:
# 	@echo "Building..."
	
	
# 	@go build -o main cmd/api/main.go

# # Run the application
# run:
# 	@go run cmd/api/main.go
# # Create DB container
# docker-run:
# 	@if docker compose up --build 2>/dev/null; then \
# 		: ; \
# 	else \
# 		echo "Falling back to Docker Compose V1"; \
# 		docker-compose up --build; \
# 	fi

# # Shutdown DB container
# docker-down:
# 	@if docker compose down 2>/dev/null; then \
# 		: ; \
# 	else \
# 		echo "Falling back to Docker Compose V1"; \
# 		docker-compose down; \
# 	fi

# # Test the application
# test:
# 	@echo "Testing..."
# 	@go test ./... -v
# # Integrations Tests for the application
# itest:
# 	@echo "Running integration tests..."
# 	@go test ./internal/database -v

# # Clean the binary
# clean:
# 	@echo "Cleaning..."
# 	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# .PHONY: all build run test clean watch docker-run docker-down itest
