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
install:
	pnpm install

# Delete Dependencies
clean:
	rm -rf pnpm-lock.yaml node_modules
	pnpm clean

# Run linting
lint:
	pnpm lint

# Formatting
format:
	pnpm format:write

format-check:
	pnpm format:check

# Type checking
tsc:
	pnpm tsc

# Line Count
line-count:
	cd $(SUBDIR) && git ls-files --others --exclude-standard --cached | grep -vE 'package-lock.json|openapi-spec.yaml|prisma/migrations/*|.pdf|.tiff' | xargs wc -l | sort -nr | awk '{print $$2, $$1}'

# Set up Minikube
minikube:
	./deploy/minikube.sh

# Seed the database
seed:
	skaffold build --file-output=build.json --profile dev --quiet
	skaffold exec seed --build-artifacts=build.json --profile dev

# Generate OpenAPI Spec and UI
generate:
	curl --fail -X GET "https://api.archesai.test/swagger/yaml"  > ./api/openapi-spec.yaml
	cd ui && npm run gen
	$(MAKE) format
	

test:
	pnpm test

test-e2e:
	skaffold build --file-output=build.json --profile dev --quiet
	skaffold exec test-e2e --build-artifacts=build.json --profile dev