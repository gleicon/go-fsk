# GO-FSK Makefile
# High-order FSK (Frequency Shift Keying) implementation

# Build configuration
BINARY_NAME = fsk-modem
MAIN_PACKAGE = ./cmd/fsk-modem
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# Go configuration
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOFMT = gofmt

# Directories
BUILD_DIR = build
EXAMPLES_DIR = examples
DIST_DIR = dist

# Default target
.PHONY: all
all: clean build examples

# Help target
.PHONY: help
help:
	@echo "GO-FSK Build System"
	@echo "==================="
	@echo ""
	@echo "Build Commands:"
	@echo "  make build          Build main fsk-modem binary"
	@echo "  make examples       Build all example programs"
	@echo "  make all            Clean and build everything"
	@echo ""
	@echo "Development Commands:"
	@echo "  make test           Run all tests"
	@echo "  make clean          Clean build artifacts"
	@echo "  make deps           Download dependencies"
	@echo "  make fmt            Format Go source code"
	@echo "  make vet            Run go vet"
	@echo "  make lint           Run static analysis"
	@echo ""
	@echo "Demo Commands:"
	@echo "  make demo           Run basic FSK demo"
	@echo "  make demo-file      Generate and decode WAV file"
	@echo "  make demo-ultrasonic Run ultrasonic example"
	@echo ""
	@echo "WebAssembly Commands:"
	@echo "  make wasm           Build WASM demo"
	@echo "  make wasm-serve     Build and serve WASM demo (http://localhost:8080)"
	@echo "  make wasm-clean     Clean WASM build artifacts"
	@echo ""
	@echo "Install Commands:"
	@echo "  make install        Install fsk-modem to \$$GOPATH/bin"
	@echo "  make uninstall      Remove fsk-modem from \$$GOPATH/bin"
	@echo ""
	@echo "Distribution:"
	@echo "  make dist           Create distribution packages"
	@echo "  make release        Tag and create release"

# Build main binary
.PHONY: build
build: deps
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

# Build all examples
.PHONY: examples
examples: deps
	@echo "Building examples..."
	@mkdir -p $(BUILD_DIR)/examples
	@for dir in $(EXAMPLES_DIR)/*/; do \
		if [ -f "$$dir/main.go" ]; then \
			example_name=$$(basename "$$dir"); \
			echo "  Building $$example_name..."; \
			$(GOBUILD) -o $(BUILD_DIR)/examples/$$example_name "$$dir/main.go"; \
		fi \
	done
	@echo "Examples built in $(BUILD_DIR)/examples/"

# Individual example builds
.PHONY: example-simple
example-simple: deps
	@echo "Building simple example..."
	@mkdir -p $(BUILD_DIR)/examples
	$(GOBUILD) -o $(BUILD_DIR)/examples/simple $(EXAMPLES_DIR)/simple/main.go

.PHONY: example-ultrasonic
example-ultrasonic: deps
	@echo "Building ultrasonic example..."
	@mkdir -p $(BUILD_DIR)/examples
	$(GOBUILD) -o $(BUILD_DIR)/examples/ultrasonic $(EXAMPLES_DIR)/ultrasonic/main.go

.PHONY: example-chat
example-chat: deps
	@echo "Building chat-tui example..."
	@mkdir -p $(BUILD_DIR)/examples
	$(GOBUILD) -o $(BUILD_DIR)/examples/chat-tui $(EXAMPLES_DIR)/chat-tui/main.go

.PHONY: example-frequency-test
example-frequency-test: deps
	@echo "Building frequency-test example..."
	@mkdir -p $(BUILD_DIR)/examples
	$(GOBUILD) -o $(BUILD_DIR)/examples/frequency-test $(EXAMPLES_DIR)/frequency-test/main.go

# Development targets
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .

.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

.PHONY: lint
lint: vet
	@echo "Running static analysis..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping linting"; \
		echo "Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Demo targets
.PHONY: demo
demo: build
	@echo "Running basic FSK test..."
	./$(BUILD_DIR)/$(BINARY_NAME) -test

.PHONY: demo-file
demo-file: build
	@echo "Generating WAV file with FSK signal..."
	./$(BUILD_DIR)/$(BINARY_NAME) -mode tx -msg "Hello FSK World!" -output demo_signal.wav
	@echo "Decoding WAV file..."
	./$(BUILD_DIR)/$(BINARY_NAME) -mode rx -input demo_signal.wav
	@echo "Generated demo_signal.wav (play with: ffplay demo_signal.wav)"

.PHONY: demo-ultrasonic
demo-ultrasonic: example-ultrasonic
	@echo "Running ultrasonic FSK demo..."
	./$(BUILD_DIR)/examples/ultrasonic

.PHONY: demo-simple
demo-simple: example-simple
	@echo "Running simple FSK demo..."
	./$(BUILD_DIR)/examples/simple

# Install targets
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to \$$GOPATH/bin..."
	$(GOCMD) install $(LDFLAGS) $(MAIN_PACKAGE)

.PHONY: uninstall
uninstall:
	@echo "Removing $(BINARY_NAME) from \$$GOPATH/bin..."
	@rm -f "$(shell go env GOPATH)/bin/$(BINARY_NAME)"

# Clean targets
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@rm -f *.wav
	@rm -f $(EXAMPLES_DIR)/*/*.wav
	@echo "Clean complete"

.PHONY: clean-examples
clean-examples:
	@echo "Cleaning example binaries..."
	@rm -rf $(BUILD_DIR)/examples

# Distribution targets
.PHONY: dist
dist: clean
	@echo "Creating distribution packages..."
	@mkdir -p $(DIST_DIR)
	
	# Linux AMD64
	@echo "Building for Linux AMD64..."
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	
	# macOS AMD64
	@echo "Building for macOS AMD64..."
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	
	# macOS ARM64 (M1/M2)
	@echo "Building for macOS ARM64..."
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	
	# Windows AMD64
	@echo "Building for Windows AMD64..."
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	
	@echo "Distribution packages created in $(DIST_DIR)/"

.PHONY: release
release:
	@echo "Creating release..."
	@if [ -z "$(TAG)" ]; then \
		echo "Usage: make release TAG=v1.0.0"; \
		exit 1; \
	fi
	@git tag -a $(TAG) -m "Release $(TAG)"
	@echo "Tagged release $(TAG)"
	@echo "Push with: git push origin $(TAG)"

# Quick development workflow
.PHONY: dev
dev: clean fmt vet build test
	@echo "Development build complete"

.PHONY: check
check: fmt vet lint test
	@echo "All checks passed"

# WASM targets
.PHONY: wasm
wasm: wasm-build wasm-copy-exec

.PHONY: wasm-build
wasm-build:
	@echo "Building WASM binary..."
	@mkdir -p wasm
	GOOS=js GOARCH=wasm go build -o wasm/fsk.wasm ./wasm/

.PHONY: wasm-copy-exec
wasm-copy-exec:
	@echo "Copying WASM executor..."
	@cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" wasm/ 2>/dev/null || \
	 cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" wasm/ 2>/dev/null || \
	 echo "Warning: Could not find wasm_exec.js, please copy manually"

.PHONY: wasm-serve
wasm-serve: wasm
	@echo "Starting development server for WASM demo..."
	@echo "Open http://localhost:8080 in your browser"
	@echo "Press Ctrl+C to stop the server"
	@cd wasm && \
	if command -v python3 >/dev/null 2>&1; then \
		python3 -m http.server 8080; \
	elif command -v python >/dev/null 2>&1; then \
		python -m SimpleHTTPServer 8080; \
	else \
		echo "Error: Python not found. Please serve the wasm/ directory manually."; \
		exit 1; \
	fi

.PHONY: wasm-clean
wasm-clean:
	@echo "Cleaning WASM build artifacts..."
	@rm -f wasm/fsk.wasm wasm/wasm_exec.js

# Docker targets (optional)
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	@if [ -f Dockerfile ]; then \
		docker build -t fsk-modem:$(VERSION) .; \
	else \
		echo "Dockerfile not found, skipping Docker build"; \
	fi

# Documentation
.PHONY: docs
docs:
	@echo "Generating documentation..."
	$(GOCMD) doc -all ./fsk

# Quick test specific examples
.PHONY: test-chat
test-chat: example-chat
	@echo "Starting chat-tui demo..."
	@echo "Press Ctrl+C to exit"
	./$(BUILD_DIR)/examples/chat-tui TestUser

.PHONY: test-frequencies
test-frequencies: example-frequency-test
	@echo "Running frequency collision test..."
	./$(BUILD_DIR)/examples/frequency-test 3

# Version info
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Build time: $(BUILD_TIME)"
	@echo "Go version: $(shell go version)"