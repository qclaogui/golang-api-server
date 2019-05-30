PROJECT?=github.com/qclaogui/golang-api-server
COMMIT?=$(shell git rev-parse --short HEAD)
RELEASE?=0.0.0
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

IMAGE?=golang-api-server
APP_PORT?=5012
APP?="main"
clean:
	rm -f ${APP}

build: clean
		go build -ldflags "-s -w \
		-X ${PROJECT}/version.Commit=${COMMIT} \
		-X ${PROJECT}/version.Release=${RELEASE} \
		-X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
		-o ${APP} cmd/main.go

container:
	docker build --build-arg APP_PORT=$(APP_PORT) \
	--build-arg COMMIT=$(COMMIT) \
	--build-arg RELEASE=$(RELEASE) \
	--build-arg BUILD_TIME=$(BUILD_TIME) \
	-t $(IMAGE):$(RELEASE) .

run: container
	docker stop $(IMAGE) || true && docker rm $(IMAGE) || true
	docker run -d --name ${IMAGE} -p ${APP_PORT}:${APP_PORT} $(IMAGE):$(RELEASE)

test:
	go test -v -race ./...