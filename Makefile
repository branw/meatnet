.PHONY: build-logger build install-logger install

build-logger:
	go build -o ./bin/meatnet-logger ./cmd/logger

build: build-logger

install-logger:
	go build -o $(shell go env GOPATH)/bin/meatnet-logger ./cmd/logger

install: install-logger
