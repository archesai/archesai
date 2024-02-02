PROFILE ?= minimal

# THESE ARE FOR RUNNING THE SERVICES IN DOCKER
docker-build:
	-make docker-stop
	-docker volume rm archesai_node_modules_ui
	-docker volume rm archesai_node_modules_api
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml build

docker-run:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) up -d

docker-db-seed:
	-docker stop arches-api
	cd api && npm run db:seed && cd ..
	docker start arches-api
	
docker-run-ui:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) up arches-ui -d

docker-test:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) up arches-api-test-e2e

docker-stop:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) down
		
docker-reset:
	make docker-build && make docker-run
	
	
# THESE ARE FOR BUILDING AND PUSHING THE SERVICES TO GCP
build-and-push-ui: 
	gcloud builds submit --config ./ui/cloudbuild.yaml --async ./ui

build-and-push-api: 
	gcloud builds submit --config ./api/cloudbuild.yaml --async ./api

build-and-push-nlp: 
	gcloud builds submit --config ./nlp/cloudbuild.yaml --async ./nlp

build-and-push-all:	
	make build-and-push-ui && make build-and-push-api && make build-and-push-nlp

# THESE ARE FOR PROXYING THE SERVICES TO LOCALHOST
telepresence-install:
	telepresence helm install

telepresence-api:
	telepresence intercept arches-api --port 3001:3000 -n archesai-dev

telepresence-ui:
	telepresence intercept arches-ui --port 3000:3000 -n archesai-dev

telepresence-all: 
	make telepresence-stop && make telepresence-api && make telepresence-ui

telepresence-stop:
	telepresence quit

REPO_LIST = ui api nlp widget

download_repos:
	@for repo in $(REPO_LIST); do \
		echo "Cloning $$repo..."; \
		git clone https://github.com/archesai/$$repo; \
	done

release:
	gcloud deploy releases create ${RELEASE} \
	--delivery-pipeline=arches-deployment \
	--region=us-central1 \
	--source=./ \
	--images=us-east4-docker.pkg.dev/archesai/images/arches-ui=us-east4-docker.pkg.dev/archesai/images/arches-ui:${UI_TAG},\
	us-east4-docker.pkg.dev/archesai/images/arches-api=us-east4-docker.pkg.dev/archesai/images/arches-api:${API_TAG},\
	us-east4-docker.pkg.dev/archesai/images/arches-nlp=us-east4-docker.pkg.dev/archesai/images/arches-nlp:${nlp_TAG}
