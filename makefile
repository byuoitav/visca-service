NAME := visca-service
OWNER := byuoitav
PKG := github.com/${OWNER}/${NAME}
DOCKER_URL := docker.pkg.github.com

# version:
# use the git tag, if this commit
# doesn't have a tag, use the git hash
COMMIT_HASH := $(shell git rev-parse --short HEAD)
TAG := $(shell git rev-parse --short HEAD)
ifneq ($(shell git describe --exact-match --tags HEAD 2> /dev/null),)
	TAG = $(shell git describe --exact-match --tags HEAD)
endif

PRD_TAG_REGEX := "v[0-9]+\.[0-9]+\.[0-9]+"
DEV_TAG_REGEX := "v[0-9]+\.[0-9]+\.[0-9]+-.+"

# go stuff
PKG_LIST := $(shell go list ${PKG}/...)

.PHONY: all deps build test test-cov clean

all: clean build

test:
	@go test -v ${PKG_LIST}

test-cov:
	@go test -coverprofile=coverage.txt -covermode=atomic ${PKG_LIST}

lint:
	@golangci-lint run --tests=false

deps:
	@echo Downloading backend dependencies...
	@go mod download

build: deps
	@mkdir -p dist

	@echo
	@echo Building for linux-amd64...
	@cd cmd/visca-service/ && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ../../dist/${NAME}-linux-amd64

	@echo
	@echo Build output is located in ./dist/.

docker: clean build
ifeq (${COMMIT_HASH}, ${TAG})
	@echo Building dev container with tag ${COMMIT_HASH}

	@echo Building container ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-dev:${COMMIT_HASH}
	@docker build -f dockerfile --build-arg NAME=${NAME}-linux-amd64 -t ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-dev:${COMMIT_HASH} dist
else ifneq ($(shell echo ${TAG} | grep -x -E ${DEV_TAG_REGEX}),)
	@echo Building dev container with tag ${TAG}

	@echo Building container ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-dev:${TAG}
	@docker build -f dockerfile --build-arg NAME=${NAME}-linux-amd64 -t ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-dev:${TAG} dist
else ifneq ($(shell echo ${TAG} | grep -x -E ${PRD_TAG_REGEX}),)
	@echo Building prd container with tag ${TAG}

	@echo Building container ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}:${TAG}
	@docker build -f dockerfile --build-arg NAME=${NAME}-linux-amd64 -t ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}:${TAG} dist
endif

deploy: docker
	@echo Logging into Github Package Registry
	@docker login ${DOCKER_URL} -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD}

ifeq (${COMMIT_HASH}, ${TAG})
	@echo Pushing dev container with tag ${COMMIT_HASH}

	@echo Pushing container ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-dev:${COMMIT_HASH}
	@docker push ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-dev:${COMMIT_HASH}
else ifneq ($(shell echo ${TAG} | grep -x -E ${DEV_TAG_REGEX}),)
	@echo Pushing dev container with tag ${TAG}

	@echo Pushing container ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-dev:${TAG}
	@docker push ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-dev:${TAG}
else ifneq ($(shell echo ${TAG} | grep -x -E ${PRD_TAG_REGEX}),)
	@echo Pushing prd container with tag ${TAG}

	@echo Pushing container ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}:${TAG}
	@docker push ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}:${TAG}
endif

clean:
	@go clean
	@rm -rf dist/
