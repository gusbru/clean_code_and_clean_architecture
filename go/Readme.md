# Go Project Setup

## Additional useful commands

`go build` - Compile the project
`go mod tidy` - Clean up dependencies
`go mod download` - Download dependencies
`go test` - Run tests

The `go.mod` file that gets created will track your dependencies and Go version requirements for the project.

## How to run the project

To run the Go application, you can use the provided Makefile commands. Here are the available commands:

### Go Commands

- `make build` - Build the Go application
- `make build-linux` - Build for Linux (cross-compilation)
- `make run` - Build and run the application
- `make dev` - Run without building a binary
- `make clean` - Clean build artifacts

### Testing Commands

- `make test` - Run tests
- `make test-nocache` - Run tests without cache
- `make test-coverage` - Run tests with coverage report
- `make test-func function=TestName` - Run specific test function

### Docker Compose Commands

- `make docker-up` - Start services with docker-compose up -d
- `make docker-down` - Stop services with docker-compose down
- `make docker-restart` - Restart docker compose services
- `make docker-exec` - Execute bash in the app container
- `make docker-shell` - Execute sh in the app container (fallback)
- `make docker-logs` - Show container logs
- `make docker-rebuild` - Rebuild and restart services
- `make docker-status` - Show services status

### Development Tools

- `make deps` - Download and tidy Go modules
- `make fmt` - Format Go code
- `make lint` - Run linting
- `make install-tools` - Install development tools (goimports, golangci-lint)

### Utility Commands

- `make help` - Show available commands
- `make all` - Run tests and build (default target)
