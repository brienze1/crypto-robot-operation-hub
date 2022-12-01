FROM golang:1.19.0-alpine

WORKDIR src/usr/crypto-robot-operation-hub

# copy source code to container
COPY . .

# build go binary
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o operation-hub cmd/operation-hub/main.go

# copy env files
RUN mkdir -p /config
COPY config/.env ./config/.env
COPY config/.env.localstack ./config/.env.localstack

# Install zip in container
RUN apk update
RUN apk add zip
RUN apk add bash

# zip the binary in the container
#RUN zip -r crypto-robot-operation-hub.zip config operation-hub

ENTRYPOINT []