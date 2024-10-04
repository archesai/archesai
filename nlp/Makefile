all:

build:
	-docker rm ml-api;
	DOCKER_BUILDKIT=1 docker build -t ml-api .;

run: build
	docker run --name ml-api -p 8080:80 ml-api;

test: build
	docker run --name ml-api -p 8080:80 ml-api python3 -m pytest;
