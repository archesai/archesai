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

# Install Dependencies
install:
	pnpm install

# Prune Dependencies
dependencies-prune:
	pnpm dedupe && pnpm prune

# Delete Dependencies
clean:
	pnpm clean:workspaces && pnpm clean

# Run linting
lint:
	pnpm lint

# Formatting
format:
	pnpm format:write

format-check:
	pnpm format:check

# Type checking
typecheck:
	pnpm typecheck

# Run valiation
validate:
	pnpm validate
	
# K8S Cluster Commands
cluster-start:
	k3d cluster create tower --config ./deploy/development/cluster-config.yaml

cluster-stop:
	k3d cluster delete -a

# Migrate the database
migrate:
	skaffold build --file-output=build.json --profile dev
	skaffold exec migrate --build-artifacts=build.json --profile dev

# Seed the database
seed:
	skaffold build --file-output=build.json --profile dev
	skaffold exec seed --build-artifacts=build.json --profile dev

generate:
	pnpm generate
	
test:
	pnpm test

test-e2e:
	skaffold build --file-output=build.json --profile dev
	skaffold exec test-e2e --build-artifacts=build.json --profile dev


# pnpm turbo build --graph --dry | dot -Tpng -oturbo-graph.png && open turbo-graph.png