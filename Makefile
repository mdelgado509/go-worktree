.PHONY: build test lint clean install fmt

# Set build variables
BINARY_NAME=go-worktree
VERSION=1.0.0
BUILD_DIR=build
INSTALL_DIR=$(HOME)/.local/bin

# Build the binary
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/go-worktree

# Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run integration tests (more invasive)
test-integration:
	@echo "Running integration tests..."
	RUN_INTEGRATION_TESTS=1 go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Checking for code issues..."
	go vet ./...

# Run code linting
lint: fmt
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping additional linting"; \
	fi

# Install the binary
install: build
	@echo "Installing to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@echo "Installation complete! Make sure $(INSTALL_DIR) is in your PATH."

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR) $(INSTALL_DIR)/$(BINARY_NAME)