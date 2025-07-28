DOCKER_DIR=./docker
APP=snake-app
APP_DATA=snake-data
APP_DATA_DIR=${DOCKER_DIR}/${APP_DATA}
APP_DIR=${DOCKER_DIR}/${APP}

model:
	@gormodel

install:
	@go install ./cmd/${APP}
	@go install ./cmd/${APP_DATA}


DOCKERFILE=${APP_DATA_DIR}/Dockerfile
docker-build-data:
	./scripts/shell/podman/build.sh ${APP_DATA} ${DOCKERFILE}
COMPOSE_FILE=${APP_DATA_DIR}/compose.yml
ENV_FILE=${APP_DATA_DIR}/.env
docker-run-data:
	./scripts/shell/podman/compose/up.sh ${APP_DATA} ${COMPOSE_FILE} ${ENV_FILE}
docker-deploy-data:
	@podman buildx build --platform linux/amd64 -f ${DOCKERFILE} -t ${APP_DATA}:local .
	./scripts/shell/podman/deploy.sh ${APP_DATA} achillesss/snake-data

DOCKERFILE_APP=${APP_DIR}/Dockerfile
docker-build-app:
	./scripts/docker/shell/local.sh ${APP} ${DOCKERFILE_APP}
COMPOSE_FILE_APP=${APP_DIR}/compose.yml
ENV_FILE_APP=${APP_DIR}/.env
docker-run-app:
	./scripts/docker/shell/compose/start.sh ${APP} ${COMPOSE_FILE_APP} ${ENV_FILE_APP} ${APP}
