# Makefile for the Arches project
PHONY_TARGETS_DEPS = install-ui-deps install-api-deps install-deps remove-ui-deps remove-api-deps remove-deps
PHONY_TARGETS_LINT_FORMAT = lint-ui lint-api lint format-ui format-api format format-check-ui format-check-api format-check line-count
PHONY_TARGETS_TEST = test-ui test-api test test-e2e
PHONY_TARGETS_DEV = run minikube seed generate
PHONY_TARGETS = $(PHONY_TARGETS_DEPS) $(PHONY_TARGETS_LINT_FORMAT) $(PHONY_TARGETS_TEST) $(PHONY_TARGETS_DEV)
.PHONY: $(PHONY_TARGETS)

# Variables
MAKEFLAGS += -j4
TEST_FILE ?= ""
SUBDIR ?= .

# Run the application
run:
	skaffold dev --profile dev --no-prune=false --cache-artifacts=false

# Install Dependencies
install-ui-deps:
	cd ui && npm install

install-api-deps:
	cd api && npm install

install-deps: install-ui-deps install-api-deps

# Clean Up and Build Files
remove-ui-deps:
	cd ui && rm -rf node_modules

remove-api-deps:
	cd api && rm -rf node_modules

remove-deps: remove-ui-deps remove-api-deps

# Linting
lint-ui:
	cd ui && npm run lint

lint-api:
	cd api && npm run lint

lint: lint-ui lint-api

# Formatting
format-ui:
	cd ui && npm run format

format-api:
	cd api && npm run format

format: format-ui format-api

format-check-ui:
	cd ui && npm run format:check

format-check-api:
	cd api && npm run format:check

format-check: format-check-ui format-check-api

# Line Count
line-count:
	cd $(SUBDIR) && git ls-files --others --exclude-standard --cached | grep -vE 'package-lock.json|openapi-spec.yaml|prisma/migrations/*|.pdf|.tiff' | xargs wc -l | sort -nr | awk '{print $$2, $$1}'

# Set up Minikube
minikube:
	./deploy/minikube.sh

# Seed the database
seed:
	skaffold build --file-output=build.json --profile dev && skaffold exec seed --build-artifacts=build.json --profile dev && rm build.json

# Generate OpenAPI Spec and UI
generate:
	curl --fail -X GET "http://arches-api.test/swagger/yaml"  > ./api/openapi-spec.yaml && cd ui && npm run gen
	$(MAKE) format
	
# Testing
test-ui:
	cd ui && npm run test

test-api:
	cd api && npm run test

test: test-ui test-api

test-e2e:
	skaffold build --file-output=build.json --profile dev && skaffold exec test-e2e --build-artifacts=build.json --profile dev && rm build.json