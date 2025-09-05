# ==========================================
# ARCHESAI MAKEFILE
# ==========================================

# ------------------------------------------
# Configuration
# ------------------------------------------

# Build Configuration
MAKEFLAGS += -j4 --no-print-directory
SERVER_OUTPUT := bin/archesai
CODEGEN_OUTPUT := bin/codegen

# Database Configuration.
MIGRATION_PATH := internal/infrastructure/database/migrations

# Terminal Colors
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
RED := \033[0;31m
NC := \033[0m # No Color

# ------------------------------------------
# Primary Targets
# ------------------------------------------

.PHONY: all
all: generate lint format ## Default: generate, lint, and format code
	@echo -e "$(GREEN)✓ All tasks complete!$(NC)"

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-25s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

# ------------------------------------------
# Development
# ------------------------------------------

.PHONY: dev
dev: ## Run all services in development mode
	@echo -e "$(YELLOW)▶ Running in development mode...$(NC)"
	@go run cmd/archesai/main.go all

.PHONY: watch
watch: ## Run with hot reload (requires air)
	@echo -e "$(YELLOW)▶ Running with hot reload...$(NC)"
	@go tool air

.PHONY: build
build: build-archesai build-codegen ## Build all binaries
	@echo -e "$(GREEN)✓ All builds complete!$(NC)"

.PHONY: build-archesai
build-archesai: ## Build archesai server binary
	@echo -e "$(YELLOW)▶ Building archesai server...$(NC)"
	@go build -o $(SERVER_OUTPUT) cmd/archesai/main.go
	@echo -e "$(GREEN)✓ archesai built: $(SERVER_OUTPUT)$(NC)"

.PHONY: build-codegen
build-codegen: ## Build codegen binary
	@echo -e "$(YELLOW)▶ Building codegen tool...$(NC)"
	@go build -o $(CODEGEN_OUTPUT) cmd/codegen/main.go
	@echo -e "$(GREEN)✓ codegen built: $(CODEGEN_OUTPUT)$(NC)"

.PHONY: run
run: run-api ## Alias for run-api

.PHONY: run-api
run-api: ## Run the API server
	@echo -e "$(YELLOW)▶ Starting API server...$(NC)"
	@go run cmd/archesai/main.go api

.PHONY: run-web
run-web: ## Run the web UI server
	@echo -e "$(YELLOW)▶ Starting web server...$(NC)"
	@go run cmd/archesai/main.go web

.PHONY: run-worker
run-worker: ## Run the background worker
	@echo -e "$(YELLOW)▶ Starting worker...$(NC)"
	@go run cmd/archesai/main.go worker

# ------------------------------------------
# Code Generation
# ------------------------------------------

.PHONY: generate
generate: generate-sqlc generate-oapi generate-defaults generate-adapters ## Generate all code
	@echo -e "$(GREEN)✓ All code generation complete!$(NC)"

.PHONY: generate-sqlc
generate-sqlc: ## Generate database code with sqlc
	@echo -e "$(YELLOW)▶ Generating sqlc code...$(NC)"
	@cd internal/infrastructure/database && go generate
	@echo -e "$(GREEN)✓ sqlc generation complete!$(NC)"

.PHONY: generate-oapi
generate-oapi: openapi-bundle ## Generate OpenAPI server code
	@echo -e "$(YELLOW)▶ Generating OpenAPI server code...$(NC)"
	@for domain in auth organizations workflows content; do \
		cd internal/domains/$$domain/generated/api && \
		{ go generate 2>&1 | grep -v "WARNING: You are using an OpenAPI 3.1.x specification" || [ $$? -eq 1 ]; } && \
		cd - > /dev/null; \
	done
	@for component in config health; do \
		cd internal/infrastructure/$$component/generated/api && \
		{ go generate 2>&1 | grep -v "WARNING: You are using an OpenAPI 3.1.x specification" || [ $$? -eq 1 ]; } && \
		cd - > /dev/null; \
	done
	@echo -e "$(GREEN)✓ OpenAPI generation complete!$(NC)"

.PHONY: generate-defaults
generate-defaults: ## Generate config defaults from OpenAPI
	@echo -e "$(YELLOW)▶ Generating config defaults...$(NC)"
	@go run cmd/codegen/main.go defaults
	@echo -e "$(GREEN)✓ Config defaults generated!$(NC)"

.PHONY: generate-adapters
generate-adapters: ## Generate type adapters between layers
	@echo -e "$(YELLOW)▶ Generating adapters...$(NC)"
	@go run cmd/codegen/main.go adapters
	@echo -e "$(GREEN)✓ Adapters generated!$(NC)"

.PHONY: generate-domain
generate-domain: ## Generate new domain scaffold (usage: make generate-domain name=billing tables=subscription,invoice)
	@echo -e "$(YELLOW)▶ Generating domain: $(name)...$(NC)"
	@go run cmd/codegen/main.go domain -name=$(name) -tables="$(tables)" -desc="$(desc)" $(if $(auth),-auth) $(if $(events),-events)
	@echo -e "$(GREEN)✓ Domain $(name) generated!$(NC)"

# ------------------------------------------
# Database
# ------------------------------------------

.PHONY: migrate
migrate: migrate-up ## Alias for migrate-up

.PHONY: migrate-up
migrate-up: ## Apply database migrations
	@echo -e "$(YELLOW)▶ Applying migrations...$(NC)"
	@go run cmd/archesai/main.go migrate up
	@echo -e "$(GREEN)✓ Migrations applied!$(NC)"

.PHONY: migrate-down
migrate-down: ## Rollback database migrations
	@echo -e "$(YELLOW)▶ Rolling back migrations...$(NC)"
	@go run cmd/archesai/main.go migrate down
	@echo -e "$(GREEN)✓ Migrations rolled back!$(NC)"

.PHONY: migrate-create
migrate-create: ## Create new migration (usage: make migrate-create name=add_users)
	@echo -e "$(YELLOW)▶ Creating migration: $(name)...$(NC)"
	@which migrate > /dev/null || (echo "Please install golang-migrate: https://github.com/golang-migrate/migrate" && exit 1)
	@go tool migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)
	@echo -e "$(GREEN)✓ Migration created!$(NC)"

# ------------------------------------------
# Testing
# ------------------------------------------

.PHONY: test
test: ## Run all tests
	@echo -e "$(YELLOW)▶ Running tests...$(NC)"
	@go test -v -cover ./...
	@echo -e "$(GREEN)✓ Tests complete!$(NC)"

.PHONY: test-coverage
test-coverage: ## Generate test coverage report
	@echo -e "$(YELLOW)▶ Generating coverage report...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo -e "$(GREEN)✓ Coverage report: coverage.html$(NC)"

# ------------------------------------------
# Code Quality
# ------------------------------------------

.PHONY: lint
lint: lint-go lint-openapi lint-node ## Run all linters
	@echo -e "$(GREEN)✓ All linting complete!$(NC)"

.PHONY: lint-go
lint-go: ## Run Go linter
	@echo -e "$(YELLOW)▶ Running Go linter...$(NC)"
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && pacman -Syu golangci-lint)
	@golangci-lint run ./...
	@echo -e "$(GREEN)✓ Go linting complete!$(NC)"

.PHONY: lint-node
lint-node: typecheck-node ## Run Node.js linter (includes typecheck)
	@echo -e "$(YELLOW)▶ Running Node.js linter...$(NC)"
	@which pnpm > /dev/null || (echo "Please install pnpm: https://pnpm.io/installation" && exit 1)
	@pnpm lint
	@echo -e "$(GREEN)✓ Node.js linting complete!$(NC)"

.PHONY: lint-openapi
lint-openapi: ## Lint OpenAPI specification
	@echo -e "$(YELLOW)▶ Linting OpenAPI spec...$(NC)"
	@pnpm --package=@redocly/cli dlx redocly --config api/redocly.yaml lint api/openapi.yaml
	@echo -e "$(GREEN)✓ OpenAPI linting complete!$(NC)"

.PHONY: typecheck-node
typecheck-node: ## Run TypeScript type checking
	@echo -e "$(YELLOW)▶ Type checking TypeScript...$(NC)"
	@which pnpm > /dev/null || (echo "Please install pnpm: https://pnpm.io/installation" && exit 1)
	@pnpm typecheck
	@echo -e "$(GREEN)✓ TypeScript type checking complete!$(NC)"

.PHONY: format
format: format-go format-node ## Format all code
	@echo -e "$(GREEN)✓ All code formatted!$(NC)"

.PHONY: format-go
format-go: ## Format Go code
	@echo -e "$(YELLOW)▶ Formatting Go code...$(NC)"
	@go fmt ./...
	@echo -e "$(GREEN)✓ Go code formatted!$(NC)"

.PHONY: format-node
format-node: ## Format Node.js/TypeScript code
	@echo -e "$(YELLOW)▶ Formatting Node.js code...$(NC)"
	@which pnpm > /dev/null || (echo "Please install pnpm: https://pnpm.io/installation" && exit 1)
	@pnpm format
	@echo -e "$(GREEN)✓ Node.js code formatted!$(NC)"

# ------------------------------------------
# OpenAPI Tools
# ------------------------------------------

.PHONY: openapi-bundle
openapi-bundle: lint-openapi ## Bundle OpenAPI into single file
	@echo -e "$(YELLOW)▶ Bundling OpenAPI spec...$(NC)"
	@pnpm --package=@redocly/cli dlx redocly --config api/redocly.yaml bundle api/openapi.yaml -o api/openapi.bundled.yaml
	@echo -e "$(GREEN)✓ OpenAPI bundled: api/openapi.bundled.yaml$(NC)"

.PHONY: openapi-split
openapi-split: lint-openapi ## Split OpenAPI into multiple files
	@echo -e "$(YELLOW)▶ Splitting OpenAPI spec...$(NC)"
	@pnpm --package=@redocly/cli dlx redocly --config api/redocly.yaml split api/openapi.bundled.yaml --outDir api/split
	@echo -e "$(GREEN)✓ OpenAPI split: api/split/$(NC)"

.PHONY: openapi-stats
openapi-stats: ## Show OpenAPI specification statistics
	@echo -e "$(YELLOW)▶ Analyzing OpenAPI spec...$(NC)"
	@pnpm --package=@redocly/cli dlx redocly --config api/redocly.yaml stats api/openapi.yaml

# ------------------------------------------
# Dependencies
# ------------------------------------------

.PHONY: deps
deps: deps-go deps-node ## Install all dependencies
	@echo -e "$(GREEN)✓ All dependencies installed!$(NC)"

.PHONY: deps-go
deps-go: ## Install Go dependencies and tools
	@echo -e "$(YELLOW)▶ Installing Go dependencies...$(NC)"
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo -e "$(GREEN)✓ Go dependencies installed!$(NC)"

.PHONY: deps-node
deps-node: ## Install Node.js dependencies
	@echo -e "$(YELLOW)▶ Installing Node.js dependencies...$(NC)"
	@which pnpm > /dev/null || (echo "Please install pnpm: https://pnpm.io/installation" && exit 1)
	@pnpm install
	@echo -e "$(GREEN)✓ Node.js dependencies installed!$(NC)"

.PHONY: deps-update
deps-update: deps-update-go deps-update-node ## Update all dependencies
	@echo -e "$(GREEN)✓ All dependencies updated!$(NC)"

.PHONY: deps-update-go
deps-update-go: ## Update Go dependencies
	@echo -e "$(YELLOW)▶ Updating Go dependencies...$(NC)"
	@go get -u ./... 2>&1 | grep -v "warning: ignoring symlink" || true
	@go mod tidy
	@echo -e "$(GREEN)✓ Go dependencies updated!$(NC)"

.PHONY: deps-update-node
deps-update-node: ## Update Node.js dependencies
	@echo -e "$(YELLOW)▶ Updating Node.js dependencies...$(NC)"
	@which pnpm > /dev/null || (echo "Please install pnpm: https://pnpm.io/installation" && exit 1)
	@pnpm update
	@echo -e "$(GREEN)✓ Node.js dependencies updated!$(NC)"

# ------------------------------------------
# Utilities
# ------------------------------------------

.PHONY: clean
clean: ## Clean build artifacts
	@echo -e "$(YELLOW)▶ Cleaning build artifacts...$(NC)"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo -e "$(GREEN)✓ Clean complete!$(NC)"

.PHONY: clean-generated
clean-generated: ## Clean all generated code
	@echo -e "$(YELLOW)▶ Cleaning generated code...$(NC)"
	@rm -rf internal/generated/api/auth/
	@rm -rf internal/generated/api/intelligence/
	@rm -rf internal/generated/api/config/
	@rm -rf internal/generated/api/health/
	@rm -rf internal/generated/api/common/
	@rm -f internal/generated/api/api.gen.go
	@rm -rf internal/generated/database/postgresql/
	@echo -e "$(GREEN)✓ Generated code cleaned!$(NC)"

.PHONY: install-tools
install-tools: ## Install development tools
	@echo -e "$(YELLOW)▶ Installing development tools...$(NC)"
	@go get -tool github.com/air-verse/air@latest
	@go get -tool github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	@echo -e "$(GREEN)✓ Tools installed!$(NC)"

.PHONY: install-completions
install-completions: ## Install shell completions guide
	@echo -e "$(BLUE)Shell Completions Installation:$(NC)"
	@echo ""
	@echo "For bash:"
	@echo "  $$ source <(archesai completion bash)"
	@echo "  $$ source <(codegen completion bash)"
	@echo ""
	@echo "For zsh:"
	@echo "  $$ source <(archesai completion zsh)"
	@echo "  $$ source <(codegen completion zsh)"
	@echo ""
	@echo -e "$(YELLOW)Add these to your shell profile to persist$(NC)"

# ------------------------------------------
# Infrastructure (Optional)
# ------------------------------------------

.PHONY: docker-run
docker-run: ## Build and run with Docker Compose
	@echo -e "$(YELLOW)▶ Starting Docker Compose...$(NC)"
	@docker-compose up --build

.PHONY: docker-stop
docker-stop: ## Stop Docker Compose services
	@echo -e "$(YELLOW)▶ Stopping Docker Compose...$(NC)"
	@docker-compose down

.PHONY: k8s-cluster-start
k8s-cluster-start: ## Start k3d cluster
	@k3d cluster create tower --config deployments/k3d/k3d.yaml

.PHONY: k8s-cluster-stop
k8s-cluster-stop: ## Stop k3d cluster
	@k3d cluster delete -a

.PHONY: k8s-deploy
k8s-deploy: ## Deploy with Helm
	@helm install dev ./helm/arches -f ./helm/dev-overrides.yaml

.PHONY: k8s-upgrade
k8s-upgrade: ## Upgrade Helm deployment
	@helm upgrade dev ./helm/arches -f ./helm/dev-overrides.yaml

.PHONY: skaffold-dev
skaffold-dev: ## Run with Skaffold in dev mode
	@skaffold dev --default-repo registry.localhost:5000 --profile dev

.PHONY: skaffold-run
skaffold-run: ## Deploy with Skaffold
	@skaffold run

.PHONY: skaffold-delete
skaffold-delete: ## Delete Skaffold deployment
	@skaffold delete --profile dev
