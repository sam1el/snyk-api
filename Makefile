.PHONY: build test lint generate clean help install-tools

# Variables
BINARY_NAME=snyk-api
BUILD_DIR=bin
GO=go
GOTEST=$(GO) test
GOVET=$(GO) vet
GOLINT=golangci-lint
ifeq ($(OS),Windows_NT)
EXT=.exe
else
EXT=
endif

help: ## Display this help message
	@echo "Snyk API Management Tool - Makefile targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

install-tools: ## Install development tools
	@echo "Installing development tools..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	@echo "Tools installed successfully!"

build: ## Build the CLI binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)$(EXT) ./cmd/snyk-api
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(EXT)"

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

test-integration: ## Run integration tests (requires SNYK_TOKEN)
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration ./...

coverage: test ## Generate test coverage report
	@echo "Generating coverage report..."
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run linters
	@echo "Running linters..."
	$(GOLINT) run

fmt: ## Format code
	@echo "Formatting code..."
	$(GO) fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

generate: ## Generate code from OpenAPI specs
	@echo "Generating API clients from OpenAPI specs..."
	$(GO) generate ./...

pull-specs: ## Pull latest OpenAPI specs
	@echo "Pulling latest API specifications..."
	@./scripts/pull-api-specs.sh

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	$(GO) mod tidy

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

install: build ## Install binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install ./cmd/snyk-api

all: clean lint test build ## Run all checks and build

dev: fmt vet lint test ## Run development checks

.DEFAULT_GOAL := help
