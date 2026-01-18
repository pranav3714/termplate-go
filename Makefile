BINARY_NAME := termplate
MODULE := $(shell go list -m)
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

LDFLAGS := -ldflags "-s -w \
	-X $(MODULE)/pkg/version.Version=$(VERSION) \
	-X $(MODULE)/pkg/version.Commit=$(COMMIT) \
	-X $(MODULE)/pkg/version.Date=$(DATE)"

BUILD_DIR := ./build/bin
COVERAGE_DIR := ./coverage

.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: setup
setup: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/evilmartians/lefthook@latest
	@if [ -f lefthook.yml ]; then lefthook install; fi

.PHONY: build
build: ## Build the binary
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

.PHONY: build-all
build-all: ## Build for all platforms
	@mkdir -p $(BUILD_DIR)
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR) $(COVERAGE_DIR) coverage.out dist/

.PHONY: test
test: ## Run unit tests
	go test -v -race -timeout 5m ./...

.PHONY: coverage
coverage: ## Generate coverage report
	@mkdir -p $(COVERAGE_DIR)
	go test -race -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1

.PHONY: fmt
fmt: ## Format code
	go fmt ./...
	goimports -w -local $(MODULE) .

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	golangci-lint run --fix ./...

.PHONY: vuln
vuln: ## Check for vulnerabilities
	govulncheck ./...

.PHONY: tidy
tidy: ## Tidy dependencies
	go mod tidy
	go mod verify

.PHONY: audit
audit: tidy vet lint vuln test ## Run all quality checks
	@echo "\033[32mAll checks passed!\033[0m"

.PHONY: ci
ci: audit coverage build ## Run full CI pipeline
	@echo "\033[32mCI completed!\033[0m"

.PHONY: release-prepare
release-prepare: ## Prepare a new release (interactive)
	@bash scripts/release.sh

.PHONY: release-dry-run
release-dry-run: ## Preview release changes without committing
	@bash scripts/release.sh --dry-run

.PHONY: release-patch
release-patch: ## Auto-increment patch version for release
	@bash scripts/release.sh --patch

.PHONY: release-dry
release-dry: ## Dry run release (GoReleaser snapshot)
	goreleaser release --snapshot --clean

.PHONY: release
release: ## Create release (GoReleaser - requires tag)
	goreleaser release --clean
