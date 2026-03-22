.PHONY: run run-dummy build swagger wire proto tidy lint

WIRE := $(shell go env GOPATH)/bin/wire
BUF  := $(shell go env GOPATH)/bin/buf

SWAG := $(shell go env GOPATH)/bin/swag

run: swagger wire
	go run ./cmd/server/...

run-dummy: swagger wire
	APP_ENV=dummy go run ./cmd/server/...

build: swagger
	go build -o bin/server ./cmd/server/...

## Generate OAS 3.1 docs from annotations (swag v2)
swagger:
	$(SWAG) init -g cmd/server/main.go -o docs --v3.1

## Install swag v2 CLI tool
swagger-install:
	go install github.com/swaggo/swag/v2/cmd/swag@latest

wire:
	$(WIRE) gen ./cmd/server/...

## Generate Go code from proto files (requires buf)
proto:
	$(BUF) generate

proto-install:
	go install github.com/bufbuild/buf/cmd/buf@latest

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

.DEFAULT_GOAL := run
