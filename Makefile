PROFILE ?= local
CONTAINER ?= arches-api

build:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml build

run:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) up -d

seed:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) up arches-api-seed

migrations:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml run --rm arches-api-seed /bin/sh -c 'npm run migrations:dev'

models:
	docker exec -it arches-ollama bash -c "echo llama3.1 mxbai-embed-large | xargs -n1 ollama pull"

lint:
	cd api && npm run lint && cd ../ui-new && npm run lint

test:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) up arches-api-test-e2e

stop:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) down
		
reset:
	-make stop
	make build && make run
	
logs:
	docker logs -f --tail=100 $(CONTAINER) 2>&1 | ccze -m ansi