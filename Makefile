APP=snake
DOCKER_DIR=./docker

model:
	@gormodel

install:
	@go install ./...

DOCKERFILE=${DOCKER_DIR}/Dockerfile
docker-build:
	./scripts/docker/shell/local.sh ${APP} ${DOCKERFILE}

COMPOSE_FILE=${DOCKER_DIR}/compose.yml
ENV_FILE=${DOCKER_DIR}/.env
docker-run:
	./scripts/docker/shell/compose/start.sh ${APP} ${COMPOSE_FILE} ${ENV_FILE} ${APP}