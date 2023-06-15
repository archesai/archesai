PROFILE ?= minimal

# THESE ARE FOR RUNNING THE SERVICES IN DOCKER
docker-build:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml build

docker-run:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) up -d

docker-run-ui:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) up arches-ui -d

docker-test:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) up arches-api-test-e2e

docker-stop:
	PROFILE=$(PROFILE) docker compose -f docker-compose.yaml -f docker-compose.dev.yaml --profile $(PROFILE) down
		
# THESE ARE FOR BUILDING AND PUSHING THE SERVICES TO GCP
build-and-push-ui: 
	gcloud builds submit --config ./ui/cloudbuild.yaml --async ./ui

build-and-push-api: 
	gcloud builds submit --config ./api/cloudbuild.yaml --async ./api

build-and-push-pyservice: 
	gcloud builds submit --config ./pyservice/cloudbuild.yaml --async ./pyservice

build-and-push-all:	
	make build-and-push-ui && make build-and-push-api && make build-and-push-pyservice

# THESE ARE FOR RUNNING THE SERVICES IN KUBERNETES
k8s-run:
	kubectl apply -f ./kubernetes/filechat -R -n filechat-dev

k8s-update:
	make k8s-run

k8s-stop:	
	kubectl delete -f ./kubernetes/filechat -R -n filechat-dev

# THESE ARE FOR PROXYING THE SERVICES TO LOCALHOST
telepresence-install:
	telepresence helm install

telepresence-api:
	telepresence intercept arches-api --port 3001:3000 -n filechat-dev

telepresence-ui:
	telepresence intercept arches-ui --port 3000:3000 -n filechat-dev

telepresence-all: 
	make telepresence-stop && make telepresence-api && make telepresence-ui

telepresence-stop:
	telepresence quit

REPO_LIST = ui api pyservice widget

download_repos:
	@for repo in $(REPO_LIST); do \
		echo "Cloning $$repo..."; \
		git clone https://github.com/filechat-io/$$repo; \
	done

release:
	gcloud deploy releases create ${RELEASE} \
	--delivery-pipeline=arches-deployment \
	--region=us-central1 \
	--source=./ \
	--images=us-west2-docker.pkg.dev/filechat-io/images/arches-ui=us-west2-docker.pkg.dev/filechat-io/images/arches-ui:${UI_TAG},\
	us-west2-docker.pkg.dev/filechat-io/images/arches-api=us-west2-docker.pkg.dev/filechat-io/images/arches-api:${API_TAG},\
	us-west2-docker.pkg.dev/filechat-io/images/arches-pyservice=us-west2-docker.pkg.dev/filechat-io/images/arches-pyservice:${PYSERVICE_TAG}
