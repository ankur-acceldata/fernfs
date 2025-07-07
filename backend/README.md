# FernFS Backend

A high-performance backend service for FernFS built with Go.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   make deps
   ```

### Running the Application

#### Development Mode
```bash
make run
```

#### Production Mode
```bash
make build
./bin/fernfs-backend
```

#### With Hot Reload
```bash
make install-dev  # Install air for hot reload
make dev
```

### API Endpoints

- `GET /` - API status
- `GET /health` - Health check

### Project Structure

```
backend/
├── cmd/
│   ├── server/          # Main server application
│   └── cli/             # CLI tools
├── internal/
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/         # Data models
│   ├── services/       # Business logic
│   └── database/       # Database layer
├── pkg/
│   ├── utils/          # Utility functions
│   └── config/         # Configuration
├── api/                # API documentation
├── scripts/            # Build and deployment scripts
├── docs/              # Documentation
└── test/              # Test files
```

### Environment Variables

- `PORT` - Server port (default: 8080)

### Development

#### Code Quality
```bash
make fmt     # Format code
make vet     # Vet code
make lint    # Lint code (requires golangci-lint)
```

#### Testing
```bash
make test              # Run tests
make test-coverage     # Run tests with coverage
```

#### Docker
```bash
make docker-build     # Build Docker image
make docker-run       # Run Docker container
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

### License

This project is licensed under the MIT License. 