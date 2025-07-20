.PHONY: docker-build docker-run compose-up compose-up-build compose-logs compose-down compose-exec start-dev test

# Docker commands
docker-build:
	docker build -t trade_platform -f docker/Dockerfile .

docker-run:
	docker run -p 3000:3000 -v $(shell pwd)/.:/app trade_platform

# Docker Compose commands
compose-up:
	docker-compose -f docker/docker-compose.yaml up -d

compose-up-build:
	docker-compose -f docker/docker-compose.yaml up -d --build

compose-logs:
	docker compose -f docker/docker-compose.yaml logs -f

compose-down:
	docker compose -f docker/docker-compose.yaml down -v

compose-exec:
	docker compose -f docker/docker-compose.yaml exec backend bash

# Development commands
start-dev:
	nodemon --exec ts-node src/api.ts

# Testing
test:
	jest --config jest.config.js

# Convenience aliases (optional)
dev: start-dev
build: docker-build
up: compose-up
down: compose-down
exec: compose-exec
logs: compose-logs