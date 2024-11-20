TEST_FILE ?= ""
SUBDIR ?= .

run:
	skaffold dev --profile dev

seed:
	skaffold build --file-output=build.json --profile dev && skaffold exec seed --build-artifacts=build.json --profile dev && rm build.json

migrations:
	cd api && DATABASE_URL="postgresql://admin:admin@localhost:5431/nestjs?schema=public" npm run db:reset && cd ..

generate:
	cd ui && npm run gen && cd .. && cd api/test && curl -X GET "http://arches-api.test/-yaml"  > openapi-spec.yaml

lint:
	cd api && npm run lint && cd ../ui && npm run lint

format:
	cd api && npm run format && cd ../ui && npm run format

line-count:
	cd $(SUBDIR) && git ls-files --others --exclude-standard --cached | grep -vE 'package-lock.json|openapi-spec.yaml|prisma/migrations/*|.pdf|.tiff' | xargs wc -l | sort -nr | awk '{print $$2, $$1}'

	
test:
	cd api && npm run test:cov && cd ..

test-e2e: generate
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) up arches-api-test-e2e

minikube:
	./deploy/scripts/minikube.sh