# Go CLI Tool Setup Guide

This document contains comprehensive research findings and step-by-step guidelines for setting up a production-quality Go CLI tool with best practices.

---

## Table of Contents

1. [Folder Structure](#1-folder-structure)
2. [CLI Framework (Cobra)](#2-cli-framework-cobra)
3. [Testing Patterns](#3-testing-patterns)
4. [Static Analysis Tools](#4-static-analysis-tools)
5. [Code Organization](#5-code-organization)
6. [Makefile Patterns](#6-makefile-patterns)
7. [Developer Experience Files](#7-developer-experience-files)
8. [Implementation Checklist](#8-implementation-checklist)

---

## 1. Folder Structure

### Recommended Layout

```
mycli/
├── main.go                       # Entry point (minimal code)
├── go.mod                        # Module definition
├── go.sum                        # Dependency checksums
├── Makefile                      # Build automation
├── .golangci.yml                 # Linter configuration
├── .editorconfig                 # Editor settings
├── .gitignore                    # Git ignore patterns
├── lefthook.yml                  # Git hooks
├── README.md                     # Project documentation
├── STRUCTURE.md                  # Folder structure docs (auto-updated)
│
├── cmd/                          # CLI commands (Cobra)
│   ├── root.go                   # Root command, global flags
│   ├── version.go                # Version command
│   └── <feature>/                # Feature subcommand groups
│       └── <command>.go
│
├── internal/                     # Private application code
│   ├── app/                      # App bootstrap/lifecycle
│   ├── config/                   # Configuration (Viper)
│   ├── handler/                  # Command handlers
│   ├── service/                  # Business logic
│   ├── repository/               # Data access layer
│   ├── model/                    # Domain models
│   ├── logger/                   # Logging (slog)
│   └── testutil/                 # Test utilities
│
├── pkg/                          # Public packages (importable)
│   ├── version/                  # Version info
│   └── validator/                # Shared validators
│
├── configs/                      # Config file templates
├── scripts/                      # Dev/CI scripts
├── docs/                         # Documentation
├── test/                         # Integration/E2E tests
│   ├── integration/
│   └── e2e/
├── build/                        # Build configs
│   ├── package/                  # Dockerfile
│   └── ci/                       # CI pipelines
└── deployments/                  # Deployment manifests
```

### Key Principles

| Directory | Rule |
|-----------|------|
| `cmd/` | Only CLI wiring, import logic from `internal/` |
| `internal/` | Private code, Go compiler enforces no external imports |
| `pkg/` | Public code, safe for external use |
| `test/` | Integration tests separate from unit tests |

### Dependency Flow

```
cmd/ → internal/handler/ → internal/service/ → internal/repository/
                ↓
         internal/model/
```

---

## 2. CLI Framework (Cobra)

### Why Cobra?

- **Industry standard**: Used by kubectl, Docker, Helm, Hugo
- **35,000+ GitHub stars**
- **Features**: Nested commands, auto help, persistent flags, shell completion
- **Integrates with Viper** for configuration

### Basic Structure

**cmd/root.go**:
```go
package cmd

import (
    "context"
    "os"
    "os/signal"
    "syscall"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
    Use:   "mycli",
    Short: "Brief description",
    Long:  `Detailed description with examples...`,
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        // Initialize logger, bind flags, etc.
        return nil
    },
    SilenceUsage:  true,  // Don't show usage on errors
    SilenceErrors: true,  // Handle errors ourselves
}

func Execute() error {
    ctx, cancel := signal.NotifyContext(
        context.Background(),
        syscall.SIGINT, syscall.SIGTERM,
    )
    defer cancel()
    return rootCmd.ExecuteContext(ctx)
}

func init() {
    cobra.OnInitialize(initConfig)
    rootCmd.PersistentFlags().StringP("config", "c", "", "config file")
    rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}

func initConfig() {
    viper.SetEnvPrefix("MYCLI")
    viper.AutomaticEnv()
    // Load config file...
}
```

**main.go**:
```go
package main

import (
    "os"
    "mycli/cmd"
)

func main() {
    if err := cmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

### Adding Commands

```go
// cmd/version.go
var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print version info",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("v1.0.0")
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}
```

### Flag Types

```go
// String flag
cmd.Flags().StringVarP(&name, "name", "n", "default", "description")

// Int flag
cmd.Flags().IntVarP(&count, "count", "c", 1, "description")

// Bool flag
cmd.Flags().BoolVar(&dryRun, "dry-run", false, "description")

// Required flag
cmd.MarkFlagRequired("name")

// Persistent flag (available to subcommands)
cmd.PersistentFlags().StringVar(&config, "config", "", "config file")
```

---

## 3. Testing Patterns

### Table-Driven Tests (Go Standard)

```go
func TestProcess(t *testing.T) {
    t.Parallel()  // Run tests in parallel

    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   "hello",
            want:    "HELLO",
            wantErr: false,
        },
        {
            name:    "empty input",
            input:   "",
            want:    "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        tt := tt  // Capture for parallel
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            got, err := Process(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("got %q, want %q", got, tt.want)
            }
        })
    }
}
```

### Testing Best Practices

1. **Use `t.Run()` for subtests** - better output, can run individually
2. **Use `t.Parallel()`** - faster test execution
3. **One assertion per test case** - clearer failure messages
4. **Test behavior, not implementation** - tests shouldn't break on refactors
5. **Colocate tests** - `foo_test.go` next to `foo.go`

### Test File Organization

```
internal/handler/
├── handler.go
├── handler_test.go      # Unit tests
├── example.go
└── example_test.go      # Unit tests

test/integration/
└── cli_test.go          # Integration tests
```

---

## 4. Static Analysis Tools

### golangci-lint (Primary Tool)

The de facto standard - runs 50+ linters in parallel with caching.

**.golangci.yml**:
```yaml
version: "2"

run:
  timeout: 5m
  modules-download-mode: readonly

linters:
  default: standard
  enable:
    # Code Quality
    - errcheck          # Unchecked errors
    - gosimple          # Simplifications
    - govet             # Suspicious constructs
    - staticcheck       # Advanced analysis
    - unused            # Unused code

    # Style
    - gofmt             # Formatting
    - goimports         # Import sorting
    - misspell          # Typos

    # Bug Detection
    - bodyclose         # HTTP body close
    - nilerr            # nil error returns
    - sqlclosecheck     # SQL close

    # Security
    - gosec             # Security issues

    # Complexity
    - cyclop            # Cyclomatic complexity
    - gocognit          # Cognitive complexity
    - nestif            # Nested ifs

    # Best Practices
    - revive            # Configurable linter
    - unconvert         # Unnecessary conversions
    - unparam           # Unused parameters
    - errorlint         # Error handling
    - wrapcheck         # Error wrapping

  settings:
    cyclop:
      max-complexity: 15
    gocognit:
      min-complexity: 20
    nestif:
      min-complexity: 5

issues:
  max-issues-per-linter: 50
  max-same-issues: 10
```

### Running Linters

```bash
# Run all linters
golangci-lint run ./...

# Auto-fix issues
golangci-lint run --fix ./...

# Show all available linters
golangci-lint help linters
```

### Other Tools

| Tool | Purpose | Command |
|------|---------|---------|
| `go vet` | Compiler-level checks | `go vet ./...` |
| `govulncheck` | Vulnerability scanning | `govulncheck ./...` |
| `staticcheck` | Advanced analysis | `staticcheck ./...` |
| `gosec` | Security scanning | `gosec ./...` |

---

## 5. Code Organization

### Logging with slog (Go 1.21+)

```go
// internal/logger/logger.go
package logger

import (
    "log/slog"
    "os"
)

func Init(level slog.Level, production bool) {
    opts := &slog.HandlerOptions{Level: level}

    var handler slog.Handler
    if production {
        handler = slog.NewJSONHandler(os.Stdout, opts)
    } else {
        handler = slog.NewTextHandler(os.Stdout, opts)
    }

    slog.SetDefault(slog.New(handler))
}

// Usage:
// slog.Info("message", "key", value)
// slog.Error("failed", "error", err)
```

### Configuration with Viper

```go
// internal/config/config.go
package config

import "github.com/spf13/viper"

type Config struct {
    Verbose  bool   `mapstructure:"verbose"`
    LogLevel string `mapstructure:"log_level"`
}

func SetDefaults() {
    viper.SetDefault("verbose", false)
    viper.SetDefault("log_level", "info")
}

func Load() (*Config, error) {
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
```

### Error Handling

```go
// DO: Handle errors at boundaries (main/CLI)
func main() {
    if err := cmd.Execute(); err != nil {
        slog.Error("failed", "error", err)
        os.Exit(1)
    }
}

// DO: Wrap errors with context
func readFile(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("reading config %s: %w", path, err)
    }
    return data, nil
}

// DON'T: Use os.Exit() in internal functions (skips defer)
// DON'T: Log and return error (choose one)
```

### Handler Pattern

```go
// internal/handler/handler.go
type Handler interface {
    Name() string
}

type ExampleInput struct {
    Name   string
    Count  int
    DryRun bool
}

type ExampleOutput struct {
    Processed int
}

// Handlers are thin - they call services
func (h *exampleHandler) Run(ctx context.Context, in ExampleInput) (*ExampleOutput, error) {
    // Validate, call service, format output
}
```

---

## 6. Makefile Patterns

```makefile
# Variables
BINARY_NAME := mycli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: help build test lint fmt audit clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

# Building
build: ## Build binary
	CGO_ENABLED=0 go build $(LDFLAGS) -o build/bin/$(BINARY_NAME) .

install: ## Install to GOPATH/bin
	go install $(LDFLAGS) .

clean: ## Remove build artifacts
	rm -rf build/bin coverage

# Testing
test: ## Run tests
	go test -v -race ./...

coverage: ## Run tests with coverage
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

bench: ## Run benchmarks
	go test -bench=. -benchmem ./...

# Code Quality
fmt: ## Format code
	go fmt ./...
	goimports -w -local $(MODULE) .

vet: ## Run go vet
	go vet ./...

lint: ## Run golangci-lint
	golangci-lint run ./...

lint-fix: ## Run golangci-lint with auto-fix
	golangci-lint run --fix ./...

vuln: ## Check vulnerabilities
	govulncheck ./...

# Combined
audit: vet lint vuln test ## Run all checks
	@echo "All checks passed!"

tidy: fmt ## Tidy and format
	go mod tidy

# Development
setup: ## Install dev tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
```

---

## 7. Developer Experience Files

### .editorconfig

```ini
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.go]
indent_style = tab
indent_size = 4

[*.{yaml,yml,json}]
indent_style = space
indent_size = 2

[Makefile]
indent_style = tab
```

### .gitignore

```gitignore
# Binaries
/build/bin/
*.exe

# Test
coverage.out
coverage.html
*.test

# IDE
.idea/
.vscode/
*.swp

# OS
.DS_Store

# Config (local)
.mycli.yaml
.env
```

### lefthook.yml (Git Hooks)

```yaml
pre-commit:
  parallel: true
  commands:
    fmt:
      glob: "*.go"
      run: gofmt -l -w {staged_files}
      stage_fixed: true
    lint:
      glob: "*.go"
      run: golangci-lint run --new-from-rev=HEAD~1

pre-push:
  commands:
    test:
      run: go test -race -short ./...
```

---

## 8. Implementation Checklist

### Phase 1: Foundation
- [ ] `go mod init`
- [ ] Create folder structure
- [ ] `main.go` (minimal)
- [ ] `cmd/root.go` with Cobra
- [ ] `.gitignore`, `.editorconfig`

### Phase 2: Core Infrastructure
- [ ] `internal/logger/` with slog
- [ ] `internal/config/` with Viper
- [ ] `pkg/version/`
- [ ] `cmd/version.go`

### Phase 3: Quality Tools
- [ ] `.golangci.yml`
- [ ] `Makefile`
- [ ] `lefthook.yml`
- [ ] `scripts/install-tools.sh`

### Phase 4: Documentation
- [ ] `README.md`
- [ ] `STRUCTURE.md`
- [ ] `configs/config.example.yaml`

### Phase 5: Example Command
- [ ] `cmd/example/` subcommand
- [ ] `internal/handler/example.go`
- [ ] Tests for all new code

### Phase 6: CI/CD (Optional)
- [ ] GitHub Actions workflow
- [ ] Dockerfile

---

## Quick Start Commands

```bash
# Initialize project
go mod init github.com/yourorg/mycli

# Install tools
make setup

# Build
make build

# Run all checks
make audit

# Run with verbose
./build/bin/mycli --verbose

# Get help
./build/bin/mycli --help
```

---

## References

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Cobra Documentation](https://cobra.dev/)
- [Viper Documentation](https://github.com/spf13/viper)
- [golangci-lint](https://golangci-lint.run/)
- [Go Wiki: TableDrivenTests](https://go.dev/wiki/TableDrivenTests)
- [Effective Go](https://go.dev/doc/effective_go)
