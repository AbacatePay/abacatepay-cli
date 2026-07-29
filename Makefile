.PHONY: build test lint clean coverage install

BINARY=abacatepay

build:
	go build -o $(BINARY) ./cmd/abacatepay

test:
	go test -v -race ./...

lint:
	golangci-lint run

install:
	go install ./cmd/abacatepay

check: test lint

deps:
	go mod download
	go mod tidy
