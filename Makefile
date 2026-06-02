.PHONY: help build test lint clean install dev version

# Build variables
VERSION ?= dev
COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)"

# Targets
BIN = tfimport-cli
BIN_PATH = ./bin/$(BIN)

help:
	@echo "tfimport-cli - Terraform State Import Block Generator"
	@echo ""
	@echo "Available targets:"
	@echo "  make build      - Build the binary for current OS"
	@echo "  make test       - Run all tests with coverage report"
	@echo "  make lint       - Run linters (vet, fmt)"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make install    - Build and install to \$$GOPATH/bin"
	@echo "  make dev        - Build and test locally"
	@echo "  make version    - Show version info"

build: clean
	@echo "Building tfimport-cli..."
	@mkdir -p bin
	@go build $(LDFLAGS) -o $(BIN_PATH) ./cmd/tfimport-cli
	@echo "✓ Binary: $(BIN_PATH)"

test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func coverage.out | tail -1
	@echo "✓ Tests passed"

lint:
	@echo "Running linters..."
	@go vet ./...
	@echo "✓ go vet passed"
	@go fmt ./...
	@echo "✓ go fmt passed"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/ dist/ *.out coverage.out
	@go clean -testcache
	@echo "✓ Clean complete"

install: build
	@echo "Installing tfimport-cli..."
	@cp $(BIN_PATH) $(GOPATH)/bin/$(BIN)
	@echo "✓ Installed to $(GOPATH)/bin/$(BIN)"

dev: lint test build
	@echo "✓ Development build complete"
	@echo "Try: $(BIN_PATH) --help"

version:
	@echo "tfimport-cli"
	@echo "  Version:   $(VERSION)"
	@echo "  Commit:    $(COMMIT)"
	@echo "  BuildDate: $(BUILD_DATE)"

.PHONY: build-all
build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BIN)-$(VERSION)-linux-amd64 ./cmd/tfimport-cli
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BIN)-$(VERSION)-linux-arm64 ./cmd/tfimport-cli
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BIN)-$(VERSION)-darwin-amd64 ./cmd/tfimport-cli
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BIN)-$(VERSION)-darwin-arm64 ./cmd/tfimport-cli
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BIN)-$(VERSION)-windows-amd64.exe ./cmd/tfimport-cli
	@echo "✓ All binaries built in dist/"

.PHONY: checksums
checksums:
	@echo "Generating checksums..."
	@cd dist && sha256sum $(BIN)-* > SHA256SUMS
	@echo "✓ Checksums: dist/SHA256SUMS"
