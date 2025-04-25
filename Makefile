APP=snake-app
APP_DATA=snake-data
DOCKER_DIR=./docker
APP_DATA_DIR=${DOCKER_DIR}/${APP_DATA}

model:
	@gormodel

install:
	@go install ./cmd/${APP}
	@go install ./cmd/${APP_DATA}


DOCKERFILE=${APP_DATA_DIR}/Dockerfile
docker-build-data:
	./scripts/docker/shell/local.sh ${APP_DATA} ${DOCKERFILE}

COMPOSE_FILE=${APP_DATA_DIR}/compose.yml
ENV_FILE=${APP_DATA_DIR}/.env
docker-run-data:
	./scripts/docker/shell/compose/start.sh ${APP_DATA} ${COMPOSE_FILE} ${ENV_FILE} ${APP_DATA}
