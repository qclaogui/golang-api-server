FROM golang:1.12.5 as builder
LABEL maintainer="qclaogui <qclaogui@gmail.com>"
ENV PROJECT github.com/qclaogui/golang-api-server
WORKDIR /root
# add source code
COPY . .
# build args, example:
# docker build --build-arg --build-arg COMMIT=$(COMMIT) --build-arg RELEASE=$(RELEASE) --build-arg BUILD_TIME=$(BUILD_TIME) -t $(IMAGE):$(RELEASE)
# commit hash
ARG COMMIT
# app build time
ARG RELEASE
# app build time
ARG BUILD_TIME
# build the source
RUN GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags "-s -w \
-X $PROJECT/version.Commit=$COMMIT \
-X $PROJECT/version.Release=$RELEASE \
-X $PROJECT/version.BuildTime=$BUILD_TIME" \
-o main cmd/main.go

# use google's best practices image
FROM gcr.io/distroless/base
# copy the binary from builder
COPY --from=builder /root/main .
# APP_PORT
ARG APP_PORT
# default 5012
ENV APP_PORT ${APP_PORT:-5012}

EXPOSE $APP_PORT
# run the binary
CMD ["./main"]