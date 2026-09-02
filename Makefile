BINARY := go-templater
CMD    := ./cmd/go-templater

.DEFAULT_GOAL := test

.PHONY: help build install run test test-v test-race cover lint vet fmt clean

help:
	@echo "Available targets:"
	@echo "  test       Run unit tests (default)"
	@echo "  test-v     Run unit tests with verbose output"
	@echo "  test-race  Run unit tests with the race detector"
	@echo "  cover      Run tests and print a per-function coverage report"
	@echo "  lint       Run golangci-lint"
	@echo "  vet        Run go vet"
	@echo "  fmt        Check formatting with gofmt"
	@echo "  build      Build the $(BINARY) binary into the repo root"
	@echo "  install    go install the $(BINARY) binary into \$$GOPATH/bin"
	@echo "  run        Run $(BINARY) (pass flags with ARGS=...)"
	@echo "  clean      Remove build and coverage artifacts"

test:
	go test ./...

test-v:
	go test -v ./...

test-race:
	go test -race ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

lint:
	golangci-lint run ./...

vet:
	go vet ./...

fmt:
	gofmt -l -s .

build:
	go build -o $(BINARY) $(CMD)

install:
	go install $(CMD)

run:
	go run $(CMD) $(ARGS)

clean:
	rm -f $(BINARY) coverage.out
