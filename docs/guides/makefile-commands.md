# Makefile Commands

Run `make help` to see all available commands.

## Available Commands

```bash
Usage: make [target]

Available targets:
  add-mapstructure-tags      Add mapstructure tags to config types for Viper compatibility
  all                        Default: generate, lint, and format code
  build                      Build all binaries
  build-api                  Build archesai server binary
  build-docs                 Build documentation site
  build-web                  Build web assets
  bundle-openapi             Bundle OpenAPI into single file
  check-deps                 Check for required dependencies
  clean                      Clean build artifacts
  clean-deps                 Clean all dependencies
  clean-generated            Clean all generated code
  clean-go                   Clean Go build artifacts
  clean-go-deps              Clean Go module cache
  clean-test                 Clean test cache and coverage files
  clean-ts                   Clean distribution builds
  clean-ts-deps              Clean Node.js dependencies
  db-migrate                 Alias for db-migrate-up
  db-migrate-create          Create new migration (usage: make db-migrate-create name=add_users)
  db-migrate-down            Rollback database migrations
  db-migrate-reset           Reset database to initial state
  db-migrate-status          Show migration status
  db-migrate-up              Apply database migrations
  deploy-docs                Manually trigger documentation deployment to GitHub Pages
  deps                       Install all dependencies
  deps-go                    Install Go dependencies and tools
  deps-ts                    Install Node.js dependencies
  deps-update                Update all dependencies
  deps-update-go             Update Go dependencies
  deps-update-ts             Update Node.js dependencies
  dev                        Run all services in development mode
  dev-all                    Run all services with hot reload
  dev-api                    Run API server with hot reload
  dev-docs                   Run documentation with hot reload
  dev-web                    Run web platform with hot reload
  docker-run                 Build and run with Docker Compose
  docker-stop                Stop Docker Compose services
  f                          Shortcut for format
  format                     Format all code
  format-go                  Format Go code
  format-prettier            Format code with Prettier
  format-ts                  Format Node.js/TypeScript code
  g                          Shortcut for generate
  generate                   Generate all code
  generate-codegen           Generate codegen
  generate-codegen-types     Generate types for codegen configuration
  generate-helm-schema       Generate Helm values.schema.json from ArchesConfig.yaml
  generate-js-client         Generate JavaScript/TypeScript client from OpenAPI
  generate-mocks             Generate test mocks using mockery
  generate-oapi              Generate OpenAPI server code
  generate-schema-sqlite     Convert PostgreSQL schema to SQLite
  generate-sqlc              Generate database code with sqlc
  help                       Show this help message
  install-tools              Install required development tools
  lint                       Run all linters
  lint-docs                  Lint documentation with markdownlint
  lint-go                    Run Go linter
  lint-openapi               Lint OpenAPI specification
  lint-ts                    Run Node.js linter (includes typecheck)
  lint-typecheck             Run TypeScript type checking
  list-workflows             List all available GitHub workflows
  prepare-docs               Copy markdown docs to web/docs/docs
  release-check              Check if ready for release
  release-clean              Clean release artifacts
  release-draft              Create a draft release on GitHub (requires gh CLI)
  release-edge-local         Test edge release workflow locally
  release-info               Show release information and next steps
  release-nightly-local      Test nightly release workflow locally
  release-snapshot           Create a snapshot release (test GoReleaser config)
  release-tag                Create and push a new release tag (usage: make release-tag VERSION=v1.0.0)
  release-test               Test release configuration without publishing
  run-api                    Run the API server (production mode)
  run-docs                   Run documentation site (production build)
  run-tui                    Launch the TUI interface
  run-web                    Run the web UI (production build)
  run-worker                 Run the background worker
  run-workflow               Run GitHub workflow locally with act (usage: make run-workflow workflow=update-docs)
  skaffold-delete            Delete Skaffold deployment
  skaffold-dev               Run with Skaffold in dev mode
  skaffold-run               Deploy with Skaffold
  split-openapi              Split OpenAPI into multiple files
  stats-openapi              Show OpenAPI specification statistics
  t                          Shortcut for test
  test                       Run all tests
  test-bench                 Run benchmark tests
  test-coverage              Generate test coverage report
  test-coverage-html         Generate HTML coverage report
  test-short                 Run short tests only (skip integration tests)
  test-verbose               Run all tests with verbose output
  test-watch                 Run tests in watch mode (requires fswatch)
  w                          Shortcut for dev-all
```
