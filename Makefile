.PHONY: build test lint install clean dev

# Build the binary
build:
	go build -o bin/cldenv cmd/cldenv/main.go

# Run tests with coverage
test:
	go test -v -race -cover ./...

# Run tests with coverage report
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run

# Install the binary
install:
	go install cmd/cldenv/main.go

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Development build (with debug info)
dev:
	go build -ldflags="-X github.com/1outres/cldenv/internal/cli.version=dev" -o bin/cldenv cmd/cldenv/main.go

# Run go mod tidy
tidy:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run vulnerability check
vuln:
	govulncheck ./...

# Run all checks
check: fmt lint test vuln

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  test          - Run tests with coverage"
	@echo "  test-coverage - Run tests with HTML coverage report"
	@echo "  lint          - Run linter"
	@echo "  install       - Install the binary"
	@echo "  clean         - Clean build artifacts"
	@echo "  dev           - Development build"
	@echo "  tidy          - Run go mod tidy"
	@echo "  fmt           - Format code"
	@echo "  vuln          - Run vulnerability check"
	@echo "  check         - Run all checks (fmt, lint, test, vuln)"
	@echo "  help          - Show this help"