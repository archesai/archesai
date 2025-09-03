# Variables
MAKEFLAGS += -j4
OUTPUT_PATH=bin/api
MAIN_PATH=cmd/api/main.go
MIGRATION_PATH=internal/database/migrations
DATABASE_URL ?= postgres://localhost/archesai?sslmode=disable

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m # No Color

.PHONY: help build run test clean clean-generated migrate generate sqlc oapi dev lint lint-go openapi-lint openapi-stats openapi-bundle fmt cluster-start cluster-stop skaffold-dev skaffold-start skaffold-stop cluster-upgrade cluster-install docker-run docker-stop deps install-tools test-coverage migrate-up migrate-down migrate-create node-deps go-deps node-update-deps go-update-deps update-deps

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${GREEN}%-15s${NC} %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development Commands
build: ## Build the application
	@echo -e "${YELLOW}Building application...${NC}"
	@go build -o $(OUTPUT_PATH) $(MAIN_PATH)
	@echo -e "${GREEN}Build complete!${NC}"

run: ## Run the application
	@echo -e "${YELLOW}Running application...${NC}"
	@go run $(MAIN_PATH)

dev: ## Run the application with hot reload (requires air)
	@echo -e "${YELLOW}Running in development mode...${NC}"
	@go tool air

# Code Generation
generate: sqlc oapi ## Generate all code (sqlc + OpenAPI)
	@echo -e "${GREEN}All code generation complete!${NC}"

sqlc: ## Generate database code with sqlc
	@echo -e "${YELLOW}Generating sqlc code...${NC}"
	@cd internal/generated/database && go generate
	@echo -e "${GREEN}sqlc generation complete!${NC}"

oapi: ## Generate OpenAPI server code
	@echo -e "${YELLOW}Generating OpenAPI server code...${NC}"
	@cd internal/generated/api && go generate
	@echo -e "${GREEN}OpenAPI generation complete!${NC}"

# Testing
test: ## Run tests
	@echo -e "${YELLOW}Running tests...${NC}"
	@go test -v -cover ./...
	@echo -e "${GREEN}Tests complete!${NC}"

test-coverage: ## Run tests with coverage report
	@echo -e "${YELLOW}Running tests with coverage...${NC}"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo -e "${GREEN}Coverage report generated: coverage.html${NC}"

# Database Migrations
migrate-up: ## Run database migrations up
	@echo -e "${YELLOW}Running migrations up...${NC}"
	@which migrate > /dev/null || (echo "Please install golang-migrate: https://github.com/golang-migrate/migrate" && exit 1)
	@go tool migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up
	@echo -e "${GREEN}Migrations complete!${NC}"

migrate-down: ## Run database migrations down
	@echo -e "${YELLOW}Running migrations down...${NC}"
	@which migrate > /dev/null || (echo "Please install golang-migrate: https://github.com/golang-migrate/migrate" && exit 1)
	@go tool migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down
	@echo -e "${GREEN}Migrations rolled back!${NC}"

migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	@echo -e "${YELLOW}Creating migration: $(name)...${NC}"
	@which migrate > /dev/null || (echo "Please install golang-migrate: https://github.com/golang-migrate/migrate" && exit 1)
	@go tool migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)
	@echo -e "${GREEN}Migration created!${NC}"

# Code Quality
lint: lint-go openapi-lint ## Run all linters (Go + OpenAPI)
	@echo -e "${GREEN}All linting complete!${NC}"

lint-go: ## Run Go linter
	@echo -e "${YELLOW}Running Go linter...${NC}"
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && pacman -Syu golangci-lint)
	@golangci-lint run ./...
	@echo -e "${GREEN}Go linting complete!${NC}"

# OpenAPI
openapi-lint: ## Lint OpenAPI specification with Redocly
	@echo -e "${YELLOW}Linting OpenAPI specification...${NC}"
	@pnpm --package=@redocly/cli dlx redocly --config api/redocly.yaml  lint api/openapi.yaml
	@echo -e "${GREEN}OpenAPI linting complete!${NC}"

openapi-stats: ## Show OpenAPI specification statistics
	@echo -e "${YELLOW}Analyzing OpenAPI specification...${NC}"
	@pnpm --package=@redocly/cli dlx redocly --config api/redocly.yaml  stats api/openapi.yaml
	@echo -e "${GREEN}OpenAPI statistics complete!${NC}"

openapi-bundle: ## Bundle OpenAPI specification into a single file
	@echo -e "${YELLOW}Bundling OpenAPI specification...${NC}"
	@pnpm --package=@redocly/cli dlx redocly --config api/redocly.yaml  bundle api/openapi.yaml -o api/openapi.bundled.yaml
	@echo -e "${GREEN}OpenAPI bundled to api/openapi.bundled.yaml${NC}"

openapi-split: ## Split OpenAPI specification into multiple files
	@echo -e "${YELLOW}Splitting OpenAPI specification...${NC}"
	@pnpm --package=@redocly/cli dlx redocly --config api/redocly.yaml  split api/openapi.bundled.yaml --outDir api
	@echo -e "${GREEN}OpenAPI split into  api/split/${NC}"

# Code Formatting
fmt: ## Format code
	@echo -e "${YELLOW}Formatting code...${NC}"
	@go fmt ./...
	@echo -e "${GREEN}Formatting complete!${NC}"

# Utilities
clean: ## Clean build artifacts
	@echo -e "${YELLOW}Cleaning...${NC}"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo -e "${GREEN}Clean complete!${NC}"

clean-generated: ## Clean all generated code
	@echo -e "${YELLOW}Cleaning generated code...${NC}"
	@rm -rf internal/generated/api/auth/
	@rm -rf internal/generated/api/intelligence/
	@rm -rf internal/generated/api/config/
	@rm -rf internal/generated/api/health/
	@rm -rf internal/generated/api/common/
	@rm -f internal/generated/api/api.gen.go
	@rm -rf internal/generated/database/postgresql/
	@echo -e "${GREEN}Generated code cleaned!${NC}"

deps: node-deps go-deps ## Download and install all dependencies	
	@echo -e "${GREEN}Dependencies updated!${NC}"

update-deps: node-update-deps go-update-deps ## Update all dependencies
	@echo -e "${GREEN}Dependencies updated!${NC}"

node-deps: ## Install Node.js dependencies
	@echo -e "${YELLOW}Installing Node.js dependencies...${NC}"
	@which pnpm > /dev/null || (echo "Please install pnpm: https://pnpm.io/installation" && exit 1)
	@pnpm install
	@echo -e "${GREEN}Node.js dependencies installed!${NC}"

node-update-deps: ## Update Node.js dependencies
	@echo -e "${YELLOW}Updating Node.js dependencies...${NC}"
	@which pnpm > /dev/null || (echo "Please install pnpm: https://pnpm.io/installation" && exit 1)
	@pnpm update
	@echo -e "${GREEN}Node.js dependencies updated!${NC}"

go-deps: ## Install Go dependencies
	@echo -e "${YELLOW}Installing Go dependencies...${NC}"
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo -e "${GREEN}Go dependencies installed!${NC}"

go-update-deps: ## Update Go dependencies
	@echo -e "${YELLOW}Updating Go dependencies...${NC}"
	@go get -u ./... 2>&1 | grep -v "warning: ignoring symlink" || true
	@go mod tidy
	@echo -e "${GREEN}Go dependencies updated!${NC}"

install-tools: ## Install development tools
	@echo -e "${YELLOW}Installing development tools...${NC}"
	@go get -tool github.com/air-verse/air@latest
	@go get -tool github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	@echo -e "${GREEN}Tools installed!${NC}"

# Kubernetes/Skaffold Commands
cluster-start: ## Start k3d cluster
	@k3d cluster create tower --config deployments/k3d/k3d.yaml

cluster-stop: ## Stop k3d cluster
	@k3d cluster delete -a

cluster-upgrade: ## Upgrade Helm deployment
	@helm upgrade dev ./helm/arches -f ./helm/dev-overrides.yaml

cluster-install: ## Install Helm deployment
	@helm install dev ./helm/arches -f ./helm/dev-overrides.yaml

skaffold-dev: ## Run application in development mode with Skaffold
	@skaffold dev --default-repo registry.localhost:5000 --profile dev

skaffold-start: ## Run application in production mode with Skaffold
	@skaffold run

skaffold-stop: ## Stop Skaffold deployment
	@skaffold delete --profile dev

# Docker
docker-run: ## Build and run the application with Docker Compose
	@docker-compose up --build

docker-stop: ## Stop Docker Compose services
	@docker-compose down