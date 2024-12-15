TEST_FILE ?= ""
SUBDIR ?= .

run:
	skaffold dev --profile dev --no-prune=false --cache-artifacts=false

seed:
	skaffold build --file-output=build.json --profile dev && skaffold exec seed --build-artifacts=build.json --profile dev && rm build.json

generate:
	curl -X GET "http://arches-api.test/swagger/yaml"  > ./api/test/openapi-spec.yaml && curl -X GET "http://arches-api.test/swagger/yaml"  > openapi-spec.yaml && cd ui && npm run gen

lint:
	cd ui && npm run lint
	cd api && npm run lint 

format:
	cd ui && npm run format
	cd api && npm run format

format-check:
	cd ui && npm run format:check
	cd api && npm run format:check

line-count:
	cd $(SUBDIR) && git ls-files --others --exclude-standard --cached | grep -vE 'package-lock.json|openapi-spec.yaml|prisma/migrations/*|.pdf|.tiff' | xargs wc -l | sort -nr | awk '{print $$2, $$1}'

test:
	cd api && npm run test:cov
	cd ui && npm run test

test-e2e: generate
	skaffold build --file-output=build.json --profile dev && skaffold exec test-e2e --build-artifacts=build.json --profile dev && rm build.json

minikube:
	./deploy/minikube.sh