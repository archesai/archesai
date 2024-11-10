PROFILE ?= stage
CONTAINER ?= arches-api
TEST_FILE ?= ""

build:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml build



build-minikube:
	eval $$(minikube docker-env) && make build

run:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) up -d

seed:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) up arches-api-seed

migrations:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml run --rm arches-api-seed /bin/sh -c 'npm run migrations:dev'

# This will curl bob:3001/-json to schema.json
generate:
	cd ui && npm run gen && cd .. && cd api/test && curl -X GET "http://bob:3001/-yaml"  > openapi-spec.yaml

models:
	docker exec -it arches-ollama bash -c "echo llama3.1 mxbai-embed-large | xargs -n1 ollama pull"

lint:
	cd api && npm run lint && cd ../ui && npm run lint

format:
	cd api && npm run format && cd ../ui && npm run format
	
test: generate
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) up arches-api-test-e2e

stop:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml down
		
reset:
	-make stop
	make build && make run
	
logs:
	docker logs -f --tail=100 $(CONTAINER) 2>&1 | ccze -m ansi