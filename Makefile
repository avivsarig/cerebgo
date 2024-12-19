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

# Creates executable - run before deploying or testing built artifact
build:
	$(GOBUILD) -v ./...

# Quick local testing of the application
run:
	$(GORUN) ./cmd/main.go

# Runs tests with race detection - use during development and CI
test:
	$(GOTEST) -v -race ./...

# Checks test coverage - run when adding new code to maintain coverage standards
coverage:
	$(GOTEST) -cover ./...
	$(GOTEST) -coverprofile=coverage.out ./...

# Catches common mistakes - run before committing code
vet:
	$(GOVET) ./...

# Comprehensive code analysis - run before committing/PR
lint:
	golangci-lint run

# Cleans up module dependencies - run after adding/removing imports
tidy:
	$(GOCMD) mod tidy

# Runs code generation - use when updating generated code (protobuf, mockgen, etc)
generate:
	$(GOGENERATE) ./...

# Tests GitHub Actions locally - use before pushing workflow changes
docker-ci:
	act push -W .github/workflows/ci.yml --container-architecture linux/amd64 # must git push first!

# Development tool - automatically runs tests when files change
watch-test:
	go install github.com/cespare/reflex@latest
	reflex -r '\.go$$' -s -- sh -c 'make test'

# Checks code formatting - use in CI to enforce standards
fmt-check:
	gofmt -l .

# Formats code - run before committing
fmt:
	gofmt -w .

# Installs development tools - run once when setting up dev environment
dev-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cespare/reflex@latest

# Removes generated files - use to ensure clean build state
clean:
	rm -f coverage.out
	rm -f $(BINARY_NAME)

# Install git hooks - run after cloning repo or updating hooks
install-hooks:
	./.scripts/install-hooks.sh

help:
	@echo "Available commands:"
	@echo "  run         		- Run the main package"
	@echo "  build       		- Build the project"
	@echo "  test        		- Run tests with race detection"
	@echo "  coverage    		- Generate test coverage report"
	@echo "  vet        		- Run go vet"
	@echo "  lint       		- Run golangci-lint"
	@echo "  tidy       		- Run go mod tidy"
	@echo "  generate   		- Run go generate"
	@echo "  docker-lint 		- Run linting in Docker"
	@echo "  docker-test 		- Run tests in Docker"
	@echo "  watch-test  		- Watch for file changes and run tests"
	@echo "  fmt-check   		- Check if files are formatted"
	@echo "  fmt        		- Format all files"
	@echo "  dev-deps   		- Install development dependencies"
	@echo "  clean      		- Remove build artifacts"
	@echo "  all        		- Run fmt, vet, lint, test, and build"
	@echo "	 install-hooks 		- Install git hooks"