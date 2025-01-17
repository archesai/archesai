# Variables
MAKEFLAGS += -j4

# Run the application in development mode
dev:
	skaffold dev --profile dev

# Run the application in production mode
start:
	skaffold run

# Stop
stop:
	skaffold delete --profile dev

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

# Run valiation
validate: lint format-check tsc
	
# K8S Cluster Commands
start-cluster:
	k3d cluster create tower --config k3d-config.yaml

stop-cluster:
	k3d cluster delete -a

# Migrate the database
migrate:
	skaffold build --file-output=build.json --profile dev
	skaffold exec migrate --build-artifacts=build.json --profile dev

# Seed the database
seed:
	skaffold build --file-output=build.json --profile dev
	skaffold exec seed --build-artifacts=build.json --profile dev

# Generate OpenAPI Spec and UI
generate:
	curl --fail -X GET "https://api.archesai.dev/swagger/yaml"  > ./api/openapi-spec.yaml
	cd ui && npm run gen
	$(MAKE) format
	
test:
	pnpm test

test-e2e:
	skaffold build --file-output=build.json --profile dev
	skaffold exec test-e2e --build-artifacts=build.json --profile dev