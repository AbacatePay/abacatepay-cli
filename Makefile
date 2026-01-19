.PHONY: build test lint clean coverage install

BINARY=abacate

build:
	go build -o $(BINARY) .

test:
	go test -v -race ./...

lint:
	golangci-lint run

install:
	go install .

check: test lint

deps:
	go mod download
	go mod tidy
