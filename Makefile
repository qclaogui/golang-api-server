PROJECT?=github.com/qclaogui/golang-api-server
COMMIT?=$(shell git rev-parse --short HEAD)
RELEASE?=0.0.0
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
PKG_LIST := $(shell go list ${PROJECT}/... | grep -v /vendor/)

IMAGE?=golang-api-server
APP_PORT?=5012
APP?="main"

.PHONY: all dep lint vet test test-coverage build clean

all: build

dep: ## Get the dependencies
	@go mod download

lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

test:
	@go test -v -race ./...

test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST}
	@cat cover.out >> coverage.txt

clean:
	@rm -f ${APP}

build: clean dep
		@go build -ldflags "-s -w \
		-X ${PROJECT}/version.Commit=${COMMIT} \
		-X ${PROJECT}/version.Release=${RELEASE} \
		-X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
		-o ${APP} cmd/main.go

#container:
#	docker build --build-arg APP_PORT=$(APP_PORT) \
#	--build-arg COMMIT=$(COMMIT) \
#	--build-arg RELEASE=$(RELEASE) \
#	--build-arg BUILD_TIME=$(BUILD_TIME) \
#	-t $(IMAGE):$(RELEASE) .
#
#run: container
#	docker stop $(IMAGE) || true && docker rm $(IMAGE) || true
#	docker run -d --name ${IMAGE} -p ${APP_PORT}:${APP_PORT} $(IMAGE):$(RELEASE)