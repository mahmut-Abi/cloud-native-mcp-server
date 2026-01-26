# Kubernetes MCP Server Makefile

# Variables
BINARY_NAME=k8s-mcp-server
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR=build
DIST_DIR=dist
MAIN_PACKAGE=./cmd/server
GO_VERSION=1.25

# Go build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -s -w"
BUILD_FLAGS=-trimpath

# Default target
.DEFAULT_GOAL := help

## help: Show this help message with detailed information
.PHONY: help
help:
	@echo
	@printf "\033[1m🚀 Kubernetes MCP Server - Makefile Help\033[0m\n"
	@echo
	@printf "\033[1mProject Information:\033[0m\n"
	@printf "  Binary Name:     $(BINARY_NAME)\n"
	@printf "  Version:         $(VERSION)\n"
	@printf "  Go Version:      $(GO_VERSION)\n"
	@printf "  Build Directory: $(BUILD_DIR)\n"
	@printf "  Dist Directory:  $(DIST_DIR)\n"
	@echo
	@printf "\033[1mUsage:\033[0m\n"
	@printf "  make \033[36m<target>\033[0m\n"
	@echo
	@printf "\033[1mQuick Start:\033[0m\n"
	@printf "  \033[36mmake build\033[0m     - Build the binary and start developing\n"
	@printf "  \033[36mmake run\033[0m       - Build and run the server locally\n"
	@printf "  \033[36mmake dev\033[0m       - Start development mode with auto-reload\n"
	@printf "  \033[36mmake check\033[0m     - Run all code quality checks\n"
	@printf "  \033[36mmake test\033[0m      - Run all tests\n"
	@echo
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo
	@printf "\033[1mCommon Workflows:\033[0m\n"
	@printf "  \033[33mDevelopment:\033[0m\n"
	@printf "    make tools && make dev\n"
	@echo
	@printf "  \033[33mTesting:\033[0m\n"
	@printf "    make check && make test-coverage\n"
	@echo
	@printf "  \033[33mRelease Preparation:\033[0m\n"
	@printf "    make clean && make build-all && make podman-build\n"
	@echo
	@printf "  \033[33mCI/CD Pipeline:\033[0m\n"
	@printf "    make deps && make check && make test-race && make build-all\n"
	@echo
	@printf "\033[1mEnvironment Variables:\033[0m\n"
	@printf "  VERSION          Override version (default: git describe)\n"
	@printf "  BUILD_DIR        Override build directory (default: build)\n"
	@printf "  DIST_DIR         Override distribution directory (default: dist)\n"
	@echo
	@printf "\033[1mRequirements:\033[0m\n"
	@printf "  - Go $(GO_VERSION)+ installed\n"
	@printf "  - Docker (for docker-* targets)\n"
	@printf "  - entr (for dev target): brew install entr / apt-get install entr\n"
	@printf "  - golangci-lint (auto-installed by make tools)\n"
	@echo
	@printf "\033[1mFor more information:\033[0m\n"
	@printf "  - Documentation: docs/\n"
	@printf "  - Configuration: config.example.yaml\n"
	@printf "  - Issues: https://github.com/mahmut-Abi/k8s-mcp-server/issues\n"
	@echo

##@ Development

## build: Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

## build-race: Build with race detector enabled
.PHONY: build-race
build-race:
	@echo "Building $(BINARY_NAME) with race detector..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 go build -race $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-race $(MAIN_PACKAGE)

## run: Run the server locally
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) --log-level debug

## dev: Run in development mode with file watching (requires entr)
.PHONY: dev
dev:
	@echo "Starting development mode..."
	@if ! command -v entr >/dev/null 2>&1; then \
		echo "entr is required for development mode. Install with: brew install entr (macOS) or apt-get install entr (Ubuntu)"; \
		exit 1; \
	fi
	find . -name "*.go" | entr -r make run

##@ Testing

## test: Run all tests with timeout
.PHONY: test
test:
	@echo "Running tests..."
	go test -timeout=2m -v ./...

## test-race: Run tests with race detector
.PHONY: test-race
test-race:
	@echo "Running tests with race detector..."
	go test -race -v ./...

## test-coverage: Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out $$(go list ./... | grep -v "/helm/handlers")
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}'); \
	THRESHOLD=$${COVERAGE_THRESHOLD:-14.0}; \
	echo "Total coverage: $$COVERAGE"; \
	echo "Minimum threshold: $$THRESHOLD%"; \
	awk -v cov="$$COVERAGE" -v thresh="$$THRESHOLD" 'BEGIN { if (cov+0 < thresh+0) { print "Error: Coverage " cov " is below minimum threshold of " thresh "%"; exit 1 } else { print "Coverage " cov " meets minimum threshold of " thresh "%" } }'

## benchmark: Run benchmarks
.PHONY: benchmark
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

##@ Code Quality

## lint: Run golangci-lint
.PHONY: lint
lint:
	@echo "Running linter..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run

## fmt: Format Go code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

## vet: Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

## check: Run all checks (fmt, vet, lint, test)
.PHONY: check
check: fmt vet lint test
	@echo "All checks passed!"

##@ Dependencies

## mod-tidy: Clean up go.mod
.PHONY: mod-tidy
mod-tidy:
	@echo "Tidying go.mod..."
	go mod tidy

## mod-verify: Verify dependencies
.PHONY: mod-verify
mod-verify:
	@echo "Verifying dependencies..."
	go mod verify

## mod-download: Download dependencies
.PHONY: mod-download
mod-download:
	@echo "Downloading dependencies..."
	go mod download

## deps: Update dependencies
.PHONY: deps
deps: mod-download mod-tidy mod-verify
	@echo "Dependencies updated!"

##@ Build & Release

## build-all: Build for all platforms
.PHONY: build-all
build-all:
	@echo "Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	
	# Linux amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	
	# Linux arm64
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)
	
	# macOS amd64
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	
	# macOS arm64
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	
	# Windows amd64
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	
	@echo "Built binaries in $(DIST_DIR)/"
	@ls -la $(DIST_DIR)/

## podman-build: Build container image using Podman
.PHONY: podman-build
podman-build:
	@echo "Building container image with Podman..."
	podman build -f deploy/Dockerfile \
		-t $(BINARY_NAME):$(VERSION) \
		-t $(BINARY_NAME):latest \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
		--build-arg VCS_REF=$(VERSION) .

## podman-run: Run container using Podman
.PHONY: podman-run
podman-run: podman-build
	@echo "Running container with Podman..."
	podman run --rm -p 8080:8080 $(BINARY_NAME):latest

## docker-build: Build container image (alias for podman-build)
.PHONY: docker-build
docker-build:
	@echo "Building container image with Podman..."
	podman build -f deploy/Dockerfile \
		-t $(BINARY_NAME):$(VERSION) \
		-t $(BINARY_NAME):latest \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
		--build-arg VCS_REF=$(VERSION) .

## docker-run: Run container (alias for podman-run)
.PHONY: docker-run
docker-run: docker-build
	@echo "Running container with Podman..."
	podman run --rm -p 8080:8080 $(BINARY_NAME):latest

##@ Cleanup

## clean: Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	rm -f coverage.out coverage.html

## clean-cache: Clean Go cache
.PHONY: clean-cache
clean-cache:
	@echo "Cleaning Go cache..."
	go clean -cache
	go clean -modcache

##@ Utilities

## install: Install binary to GOPATH/bin
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) $(MAIN_PACKAGE)

## version: Show version information
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Go version: $(shell go version)"
	@echo "Git commit: $(shell git rev-parse HEAD 2>/dev/null || echo 'unknown')"

## tools: Install development tools
.PHONY: tools
tools:
	@echo "Installing development tools..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

## e2e: Run end-to-end tests
.PHONY: e2e
e2e: build
	@echo "Running end-to-end tests..."
	@if [ -f scripts/e2e_full.sh ]; then \
		chmod +x scripts/e2e_full.sh && ./scripts/e2e_full.sh; \
	else \
		echo "e2e test script not found"; \
	fi

## demo: Run demo script
.PHONY: demo
demo: build
	@echo "Running demo..."
	@if [ -f demo.sh ]; then \
		chmod +x demo.sh && ./demo.sh; \
	else \
		echo "demo.sh not found"; \
	fi