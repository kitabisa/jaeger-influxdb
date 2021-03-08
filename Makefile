GOOS                 ?= $(shell go env GOOS)
SHELL                 = /bin/bash

APP_NAME              = jaeger-influxdb
VERSION               = $(shell git describe --always --tags)
GIT_COMMIT            = $(shell git rev-parse HEAD)
GIT_DIRTY             = $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE            = $(shell date '+%Y-%m-%d-%H:%M:%S')
SQUAD                 = infra
BUSINESS              = engineering

.PHONY: build
build:
	CGO_ENABLED=0 go build -o ./bin/jaeger-influxdb/jaeger-influxdb-$(GOOS) ./cmd/jaeger-influxdb/main.go

.PHONY: build-linux
build-linux:
	GOOS=linux $(MAKE) build

.PHONY: build-docker-local
build-docker-local: build-linux
	docker build -f Dockerfile.local -t jaeger-influx-grpc:local .

.PHONY: package
package:
	@echo "Build, tag, and push Docker image ${APP_NAME} ${VERSION} ${GIT_COMMIT}"
	docker buildx build \
		--build-arg VERSION=${VERSION},GIT_COMMIT=${GIT_COMMIT}${GIT_DIRTY} \
		--cache-from type=local,src=/tmp/.buildx-cache \
		--cache-to type=local,dest=/tmp/.buildx-cache \
		--tag ${DOCKER_REPOSITORY}/${APP_NAME}:${GIT_COMMIT} \
		--tag ${DOCKER_REPOSITORY}/${APP_NAME}:${VERSION} \
		--tag ${DOCKER_REPOSITORY}/${APP_NAME}:latest \
		--push .

.PHONY: deploy
deploy:
	@echo "Deploying ${APP_NAME} ${VERSION}"
	export APP_NAME=${APP_NAME} && \
	export VERSION=${VERSION} && \
	export SQUAD=${SQUAD} && \
	export BUSINESS=${BUSINESS} && \
	helmfile apply

.PHONY: helm-history-length
helm-history-length:
	@helm history \
		--namespace ${APP_NAME} \
		--output yaml \
		${APP_NAME}-server-${ENV} | yq r - --length

.PHONY: helm-oldest-revision
helm-oldest-revision:
	@helm history \
		--namespace ${APP_NAME} \
		--output yaml \
		${APP_NAME}-server-${ENV} | yq r - "[0]".revision

.PHONY: helm-image-tag
helm-image-tag:
	@helm get values \
		--namespace ${APP_NAME} \
		--revision ${REVISION} \
		--output yaml \
		${APP_NAME}-server-${ENV} | yq r - image.tag

.PHONY: prune
prune:
	@echo "Removing Docker image ${DOCKER_REPOSITORY}/${APP_NAME}:${IMAGE_TAG}"
	gcloud container images delete \
		--force-delete-tags \
		--quiet \
		${DOCKER_REPOSITORY}/${APP_NAME}:${IMAGE_TAG}

.PHONY: rollback
rollback:
	@echo "Rollback ${RELEASE} ${REVISION}"
	helm rollback \
		--namespace ${APP_NAME} \
		${RELEASE} ${REVISION}

.PHONY: clean
clean:
	@echo "Removing ${APP_NAME} ${VERSION}"
	@test ! -e bin/${APP_NAME} || rm bin/${APP_NAME}
