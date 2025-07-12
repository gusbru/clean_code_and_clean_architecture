# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTESTFOLDER=./tests
GOTESTDIR=$(GOTESTFOLDER)/...
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=app
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker parameters
DOCKER_COMPOSE=docker-compose
DOCKER_COMPOSE_FILE=docker/docker-compose.yaml
SERVICE_NAME=app

.PHONY: all build clean test deps run docker-up docker-down docker-exec docker-logs help

# Default target
all: test build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...

# Build for linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
test:
	$(GOTEST) -v $(GOTESTDIR)

# Run tests without cache
test-nocache:
	$(GOTEST) -v -count=1 $(GOTESTDIR)

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out $(GOTESTDIR)
	$(GOCMD) tool cover -html=coverage.out

# Test specific function. Take the function name as an argument
test-func:
	@if [ -z "$(function)" ]; then \
        echo "Usage: make test-func function=TestName"; \
        exit 1; \
    fi
	$(GOTEST) -v -count=1 -run $(function) $(GOTESTFOLDER)

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

# Run without building binary
dev:
	$(GOCMD) run ./...

# Docker Compose commands
docker-up:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d

docker-down:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down

docker-restart: docker-down docker-up

docker-logs:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) logs -f

# Execute bash in the app container
docker-exec:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) exec $(SERVICE_NAME) /bin/bash

# Execute bash in the app container (if bash is not available, try sh)
docker-shell:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) exec $(SERVICE_NAME) /bin/sh

# Rebuild and restart docker services
docker-rebuild:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) build --no-cache
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d

# View docker compose services status
docker-status:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) ps

# Install tools
install-tools:
	$(GOCMD) install golang.org/x/tools/cmd/goimports@latest
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Format code
fmt:
	$(GOCMD) fmt ./...
	goimports -w .

# Lint code
lint:
	golangci-lint run

# Help target
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  build-linux   - Build for Linux"
	@echo "  clean         - Clean build files"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  deps          - Download dependencies"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run without building binary"
	@echo "  docker-up     - Start docker compose services"
	@echo "  docker-down   - Stop docker compose services"
	@echo "  docker-restart- Restart docker compose services"
	@echo "  docker-exec   - Execute bash in app container"
	@echo "  docker-shell  - Execute sh in app container"
	@echo "  docker-logs   - Show docker compose logs"
	@echo "  docker-rebuild- Rebuild and restart services"
	@echo "  docker-status - Show services status"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  install-tools - Install development tools"
	@echo "  help          - Show this help"
