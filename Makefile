.PHONY: build run test coverage lint vet tidy generate docker-lint docker-test all clean fmt fmt-check watch-test dev-deps

# Default target
all: fmt vet lint test build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOGENERATE=$(GOCMD) generate
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
BINARY_NAME=task-manager

build:
	$(GOBUILD) -v ./...

run:
	$(GORUN) ./cmd/main.go

test:
	$(GOTEST) -v -race ./...

coverage:
	$(GOTEST) -cover ./...
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

vet:
	$(GOVET) ./...

lint:
	golangci-lint run -v

tidy:
	$(GOCMD) mod tidy

generate:
	$(GOGENERATE) ./...

docker-lint:
	act push -W .github/workflows/lint.yml --container-architecture linux/amd64

docker-test:
	act push -W .github/workflows/test.yml --container-architecture linux/amd64

watch-test:
	go install github.com/cespare/reflex@latest
	reflex -r '\.go$$' -s -- sh -c 'make test'

fmt-check:
	gofmt -l .

fmt:
	gofmt -w .

dev-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cespare/reflex@latest
	
clean:
	rm -f coverage.out
	rm -f $(BINARY_NAME)

help:
	@echo "Available commands:"
	@echo "  build       - Build the project"
	@echo "  run         - Run the main package"
	@echo "  test        - Run tests with race detection"
	@echo "  coverage    - Generate test coverage report"
	@echo "  vet        - Run go vet"
	@echo "  lint       - Run golangci-lint"
	@echo "  tidy       - Run go mod tidy"
	@echo "  generate   - Run go generate"
	@echo "  docker-lint - Run linting in Docker"
	@echo "  docker-test - Run tests in Docker"
	@echo "  watch-test  - Watch for file changes and run tests"
	@echo "  fmt-check   - Check if files are formatted"
	@echo "  fmt        - Format all files"
	@echo "  dev-deps   - Install development dependencies"
	@echo "  clean      - Remove build artifacts"
	@echo "  all        - Run fmt, vet, lint, test, and build"