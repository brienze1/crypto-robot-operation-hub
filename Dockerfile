FROM golang:1.19.0-alpine

WORKDIR /usr/src

# copy source code to container
COPY . .

# build go binary
RUN go mod download
RUN GOPATH=$GOPATH:/root go build -o crypto-robot-operation-hub/operation-hub cmd/operation-hub/main.go

# copy env files
RUN mkdir -p /usr/src/crypto-robot-operation-hub/config
COPY config/.env ./crypto-robot-operation-hub/config
COPY config/.env.localstack ./crypto-robot-operation-hub/config

# Install zip in container
RUN apk update
RUN apk add zip
RUN apk add bash

# zip the binary in the container
RUN zip -r crypto-robot-operation-hub.zip crypto-robot-operation-hub

ENTRYPOINT []