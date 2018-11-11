FROM golang:1.11.2 as builder
MAINTAINER qclaogui <qclaogui@gmail.com>

WORKDIR /go/src/github.com/qclaogui/golang-api-server
# add source code
COPY . .
# build the source
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

# use a minimal alpine image
FROM alpine:3.8
# add ca-certificates in case you need them
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
# set working directory
WORKDIR /root
# set app port
ENV APP_PORT 5012
# copy the binary from builder
COPY --from=builder /go/src/github.com/qclaogui/golang-api-server .

EXPOSE $APP_PORT
# run the binary
CMD ["./app"]