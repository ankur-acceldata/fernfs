.PHONY: build run test clean deps fmt vet lint docker-build docker-run

# Go parameters
BINARY_NAME=fernfs-server
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./cmd/server
GO=go

# Build the binary
build:
	$(GO) build -o $(BINARY_PATH) $(MAIN_PATH)

# Run the application
run:
	$(GO) run $(MAIN_PATH)

# Run the application with hot reload (requires air)
dev:
	air -c .air.toml

# Test all packages
test:
	$(GO) test -v ./...

# Test with coverage
test-coverage:
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

# Clean build artifacts
clean:
	$(GO) clean
	rm -rf ./bin/
	rm -rf ./tmp/

# Download dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

# Format code
fmt:
	$(GO) fmt ./...

# Vet code
vet:
	$(GO) vet ./...

# Lint code (requires golangci-lint)
lint:
	$(GO) vet ./...

# Install development dependencies
install-dev:
	$(GO) install github.com/cosmtrek/air@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Build Docker image
docker-build:
	docker build -t $(BINARY_NAME) .

# Run Docker container
docker-run:
	docker run -p 8080:8080 $(BINARY_NAME)

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  run           - Run the application"
	@echo "  dev           - Run with hot reload"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  lint          - Lint code"
	@echo "  install-dev   - Install development dependencies"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container" 