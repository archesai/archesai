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
