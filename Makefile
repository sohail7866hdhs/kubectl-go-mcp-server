# kubectl-go-mcp-server Makefile

# Build variables
BINARY_NAME=kubectl-go-mcp-server
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Build targets
.PHONY: all build build-all clean test test-coverage test-verbose deps fmt vet lint release release-snapshot release-local docker-build install run help

all: clean deps fmt vet test build

## build: Build the binary
build:
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd

## build-all: Build for multiple platforms
build-all: clean deps fmt vet test
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 ./cmd
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd

## clean: Remove build artifacts
clean:
	$(GOCLEAN)
	rm -rf bin/

## test: Run tests
test:
	$(GOTEST) -v ./test/...

## test-coverage: Run tests with coverage report
test-coverage:
	$(GOTEST) -coverpkg=kubectl-go-mcp-server/internal/cli,kubectl-go-mcp-server/internal/config,kubectl-go-mcp-server/internal/mcp,kubectl-go-mcp-server/pkg/kubectl,kubectl-go-mcp-server/pkg/types -coverprofile=coverage.out ./test/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@$(GOCMD) tool cover -func=coverage.out | tail -1

## test-verbose: Run tests with verbose output and coverage
test-verbose:
	$(GOTEST) -v -coverpkg=kubectl-go-mcp-server/internal/cli,kubectl-go-mcp-server/internal/config,kubectl-go-mcp-server/internal/mcp,kubectl-go-mcp-server/pkg/kubectl,kubectl-go-mcp-server/pkg/types ./test/...

## deps: Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## fmt: Format Go code
fmt:
	$(GOFMT) -s -w .

## vet: Run go vet
vet:
	$(GOCMD) vet ./...

## lint: Run golangci-lint (requires golangci-lint to be installed)
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install it from https://golangci-lint.run/welcome/install/" && exit 1)
	golangci-lint run

## run: Run the application
run:
	$(GOCMD) run ./cmd

## install: Install the binary to GOPATH/bin
install:
	$(GOCMD) install $(LDFLAGS) ./cmd

## mod-update: Update all dependencies
mod-update:
	$(GOMOD) get -u all
	$(GOMOD) tidy

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test

## docker-build: Build Docker image
docker-build:
	docker build -t $(BINARY_NAME):$(VERSION) .

## release: Create a new release tag
release:
	@echo "Creating release v$(VERSION)"
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)
	@echo "Release v$(VERSION) has been created and pushed. GitHub Actions will now build and publish the release."

## release-snapshot: Test the release process with a snapshot build
release-snapshot:
	@which goreleaser > /dev/null || (echo "goreleaser not found. Install it from https://goreleaser.com/install/" && exit 1)
	goreleaser release --snapshot --clean --config .goreleaser.local.yml

## release-local: Create a full release build locally (without publishing)
release-local:
	@which goreleaser > /dev/null || (echo "goreleaser not found. Install it from https://goreleaser.com/install/" && exit 1)
	goreleaser release --clean --skip=publish --config .goreleaser.local.yml

## help: Show this help message
help:
	@echo "kubectl-go-mcp-server Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
