# Go CLI Tool - Comprehensive Reference

> **Document Type**: Complete Reference
> **Purpose**: Authoritative source for all Go CLI patterns and best practices
> **Keywords**: go cli, cobra, viper, patterns, reference, best practices, project structure, testing, ci/cd
> **Related**: CONVENTIONS.md, GETTING_STARTED.md

A complete reference for building production-quality Go CLI tools. Covers everything from project structure to CI/CD, profiling, and release automation.

**This document is the source of truth for all Go patterns used in this project.**

---

## Table of Contents

1. [Folder Structure](#1-folder-structure)
2. [CLI Framework (Cobra)](#2-cli-framework-cobra)
3. [Configuration (Viper)](#3-configuration-viper)
4. [Logging (slog)](#4-logging-slog)
5. [Error Handling](#5-error-handling)
6. [Testing Patterns](#6-testing-patterns)
7. [Static Analysis Tools](#7-static-analysis-tools)
8. [Makefile](#8-makefile)
9. [Developer Experience Files](#9-developer-experience-files)
10. [CI/CD (GitHub Actions)](#10-cicd-github-actions)
11. [Docker](#11-docker)
12. [Releasing (GoReleaser)](#12-releasing-goreleaser)
13. [Shell Completion](#13-shell-completion)
14. [Profiling & Benchmarking](#14-profiling--benchmarking)
15. [Implementation Checklist](#15-implementation-checklist)
16. [Code Templates](#16-code-templates)

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
├── .goreleaser.yml               # Release configuration
├── .editorconfig                 # Editor settings
├── .gitignore                    # Git ignore patterns
├── lefthook.yml                  # Git hooks
├── README.md                     # Project documentation
├── STRUCTURE.md                  # Folder structure docs (auto-updated)
├── LICENSE                       # License file
├── CHANGELOG.md                  # Version history
│
├── cmd/                          # CLI commands (Cobra)
│   ├── root.go                   # Root command, global flags
│   ├── version.go                # Version command
│   ├── completion.go             # Shell completion command
│   └── <feature>/                # Feature subcommand groups
│       ├── <feature>.go          # Parent command
│       └── <subcommand>.go       # Subcommands
│
├── internal/                     # Private application code
│   ├── app/                      # App bootstrap/lifecycle
│   │   ├── app.go
│   │   └── app_test.go
│   ├── config/                   # Configuration (Viper)
│   │   ├── config.go
│   │   ├── config_test.go
│   │   └── defaults.go
│   ├── handler/                  # Command handlers (thin layer)
│   │   ├── handler.go            # Interfaces
│   │   ├── <feature>.go          # Implementation
│   │   └── <feature>_test.go
│   ├── service/                  # Business logic
│   │   └── <feature>/
│   │       ├── service.go
│   │       └── service_test.go
│   ├── repository/               # Data access layer
│   │   ├── repository.go         # Interfaces
│   │   └── <impl>/               # Implementations (memory, file, db)
│   │       ├── repository.go
│   │       └── repository_test.go
│   ├── model/                    # Domain models
│   │   ├── model.go
│   │   └── errors.go             # Domain errors
│   ├── logger/                   # Logging (slog)
│   │   ├── logger.go
│   │   └── logger_test.go
│   ├── validator/                # Input validation
│   │   ├── validator.go
│   │   └── validator_test.go
│   └── testutil/                 # Test utilities
│       ├── fixtures.go           # Test data
│       ├── mocks.go              # Mock implementations
│       └── helpers.go            # Test helpers
│
├── pkg/                          # Public packages (importable by others)
│   ├── version/                  # Version info
│   │   ├── version.go
│   │   └── version_test.go
│   └── <shared>/                 # Other shared utilities
│
├── api/                          # API definitions
│   ├── openapi/                  # OpenAPI/Swagger specs
│   └── proto/                    # Protocol Buffer definitions
│
├── configs/                      # Configuration templates
│   ├── config.example.yaml       # Example config
│   ├── config.development.yaml   # Dev defaults
│   └── config.production.yaml    # Prod defaults
│
├── scripts/                      # Development scripts
│   ├── install-tools.sh          # Install dev dependencies
│   ├── generate-structure.sh     # Generate STRUCTURE.md
│   └── coverage.sh               # Coverage report
│
├── docs/                         # Documentation
│   ├── architecture.md           # Architecture decisions
│   ├── contributing.md           # Contribution guide
│   └── commands/                 # Command reference
│       └── README.md
│
├── test/                         # Integration & E2E tests
│   ├── integration/              # Integration tests
│   │   └── cli_test.go
│   ├── e2e/                      # End-to-end tests
│   │   └── workflow_test.go
│   └── testdata/                 # Test fixtures
│       └── fixtures/
│
├── build/                        # Build configurations
│   ├── package/                  # Packaging
│   │   ├── Dockerfile
│   │   ├── Dockerfile.alpine
│   │   └── docker-compose.yml
│   └── ci/                       # CI configs
│       └── .github/
│           └── workflows/
│               ├── ci.yml
│               └── release.yml
│
├── deployments/                  # Deployment configurations
│   ├── kubernetes/               # K8s manifests
│   └── terraform/                # Infrastructure as code
│
└── tools/                        # Tool dependencies
    └── tools.go                  # Pin tool versions
```

### Key Principles

| Directory | Purpose | Rule |
|-----------|---------|------|
| `cmd/` | CLI commands | Only wiring, no business logic |
| `internal/` | Private code | Cannot be imported externally |
| `pkg/` | Public packages | Stable API, versioned |
| `test/` | Integration tests | Separate from unit tests |

### Dependency Flow

```
main.go
    ↓
cmd/ (CLI wiring)
    ↓
internal/handler/ (thin layer, validation)
    ↓
internal/service/ (business logic)
    ↓
internal/repository/ (data access)
    ↓
internal/model/ (domain entities)
```

**Rules:**
- Dependencies only flow downward
- No circular imports
- `model/` has no dependencies on other internal packages

---

## 2. CLI Framework (Cobra)

### Why Cobra?

- Industry standard (kubectl, Docker, Helm, Hugo, GitHub CLI)
- 35,000+ GitHub stars
- Features: nested commands, auto help, persistent flags, shell completion
- Integrates with Viper for configuration

### Installation

```bash
go get -u github.com/spf13/cobra@latest
```

### main.go (Minimal)

```go
package main

import (
    "os"

    "github.com/yourorg/mycli/cmd"
)

func main() {
    if err := cmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

### cmd/root.go (Full Example)

```go
package cmd

import (
    "context"
    "fmt"
    "log/slog"
    "os"
    "os/signal"
    "syscall"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "github.com/yourorg/mycli/internal/config"
    "github.com/yourorg/mycli/internal/logger"
)

var (
    cfgFile string
    verbose bool
    output  string
)

var rootCmd = &cobra.Command{
    Use:   "mycli",
    Short: "A brief description of your CLI",
    Long: `A longer description that spans multiple lines.

Examples:
  mycli --help
  mycli version
  mycli do-something --flag value`,

    // Runs before any subcommand
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        // Skip for completion and help
        if cmd.Name() == "completion" || cmd.Name() == "help" {
            return nil
        }

        // Initialize logger
        level := slog.LevelInfo
        if verbose {
            level = slog.LevelDebug
        }
        logger.Init(level, os.Getenv("ENV") == "production")

        // Bind flags to viper
        if err := viper.BindPFlags(cmd.Flags()); err != nil {
            return fmt.Errorf("binding flags: %w", err)
        }

        return nil
    },

    SilenceUsage:  true,  // Don't show usage on error
    SilenceErrors: true,  // We handle errors ourselves
}

// Execute is the entry point called from main
func Execute() error {
    // Set up context with signal handling
    ctx, cancel := signal.NotifyContext(
        context.Background(),
        syscall.SIGINT,
        syscall.SIGTERM,
    )
    defer cancel()

    return rootCmd.ExecuteContext(ctx)
}

func init() {
    cobra.OnInitialize(initConfig)

    // Persistent flags (available to all subcommands)
    rootCmd.PersistentFlags().StringVarP(
        &cfgFile,
        "config", "c",
        "",
        "config file (default: $HOME/.mycli.yaml)",
    )
    rootCmd.PersistentFlags().BoolVarP(
        &verbose,
        "verbose", "v",
        false,
        "enable verbose output",
    )
    rootCmd.PersistentFlags().StringVarP(
        &output,
        "output", "o",
        "text",
        "output format (text, json, yaml)",
    )

    // Add subcommands
    rootCmd.AddCommand(versionCmd)
    rootCmd.AddCommand(completionCmd)
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        if err != nil {
            slog.Error("failed to get home directory", "error", err)
            return
        }

        // Search paths
        viper.AddConfigPath(home)
        viper.AddConfigPath(".")
        viper.SetConfigType("yaml")
        viper.SetConfigName(".mycli")
    }

    // Environment variables
    viper.SetEnvPrefix("MYCLI")
    viper.AutomaticEnv()

    // Set defaults
    config.SetDefaults()

    // Read config (ignore if not found)
    if err := viper.ReadInConfig(); err == nil {
        slog.Debug("using config file", "file", viper.ConfigFileUsed())
    }
}
```

### cmd/version.go

```go
package cmd

import (
    "encoding/json"
    "fmt"

    "github.com/spf13/cobra"
    "gopkg.in/yaml.v3"

    "github.com/yourorg/mycli/pkg/version"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print version information",
    Long:  `Print the version, commit, build date, and Go version.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        info := version.Get()

        switch output {
        case "json":
            data, err := json.MarshalIndent(info, "", "  ")
            if err != nil {
                return err
            }
            fmt.Println(string(data))
        case "yaml":
            data, err := yaml.Marshal(info)
            if err != nil {
                return err
            }
            fmt.Print(string(data))
        default:
            fmt.Printf("mycli %s\n", info.String())
        }

        return nil
    },
}
```

### Subcommand Group Example

**cmd/project/project.go**:
```go
package project

import "github.com/spf13/cobra"

// Cmd is the parent command for project operations
var Cmd = &cobra.Command{
    Use:   "project",
    Short: "Manage projects",
    Long:  `Create, list, and manage projects.`,
}

func init() {
    Cmd.AddCommand(createCmd)
    Cmd.AddCommand(listCmd)
    Cmd.AddCommand(deleteCmd)
}
```

**cmd/project/create.go**:
```go
package project

import (
    "context"
    "fmt"
    "log/slog"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "github.com/yourorg/mycli/internal/handler"
)

var (
    name        string
    description string
    template    string
    dryRun      bool
)

var createCmd = &cobra.Command{
    Use:   "create",
    Short: "Create a new project",
    Long: `Create a new project with the specified name and template.

Examples:
  mycli project create --name myproject
  mycli project create --name myproject --template go-api
  mycli project create --name myproject --dry-run`,

    Args: cobra.NoArgs,

    PreRunE: func(cmd *cobra.Command, args []string) error {
        if name == "" {
            return fmt.Errorf("--name is required")
        }
        return nil
    },

    RunE: func(cmd *cobra.Command, args []string) error {
        return runCreate(cmd.Context())
    },
}

func init() {
    createCmd.Flags().StringVarP(&name, "name", "n", "", "project name (required)")
    createCmd.Flags().StringVarP(&description, "description", "d", "", "project description")
    createCmd.Flags().StringVarP(&template, "template", "t", "default", "project template")
    createCmd.Flags().BoolVar(&dryRun, "dry-run", false, "simulate without making changes")

    _ = createCmd.MarkFlagRequired("name")
    _ = viper.BindPFlag("project.template", createCmd.Flags().Lookup("template"))
}

func runCreate(ctx context.Context) error {
    slog.Debug("creating project",
        "name", name,
        "template", template,
        "dry_run", dryRun,
    )

    h := handler.NewProjectHandler()
    result, err := h.Create(ctx, handler.ProjectCreateInput{
        Name:        name,
        Description: description,
        Template:    template,
        DryRun:      dryRun,
    })
    if err != nil {
        return fmt.Errorf("creating project: %w", err)
    }

    if dryRun {
        fmt.Println("Dry run - no changes made")
        return nil
    }

    fmt.Printf("Created project %q at %s\n", result.Name, result.Path)
    return nil
}
```

### Flag Types Reference

```go
// String
cmd.Flags().StringVarP(&myVar, "flag", "f", "default", "description")

// Int
cmd.Flags().IntVarP(&count, "count", "c", 10, "number of items")

// Bool
cmd.Flags().BoolVar(&force, "force", false, "force operation")

// String slice
cmd.Flags().StringSliceVarP(&tags, "tag", "t", []string{}, "tags (can be repeated)")

// Duration
cmd.Flags().DurationVar(&timeout, "timeout", 30*time.Second, "timeout duration")

// Count (increments each use: -v -v -v)
cmd.Flags().CountVarP(&verbosity, "verbose", "v", "increase verbosity")

// Required flag
cmd.MarkFlagRequired("name")

// Mutually exclusive flags
cmd.MarkFlagsMutuallyExclusive("json", "yaml")

// Flags that must be used together
cmd.MarkFlagsRequiredTogether("username", "password")

// Hidden flag
cmd.Flags().BoolVar(&debug, "debug", false, "debug mode")
cmd.Flags().MarkHidden("debug")

// Deprecated flag
cmd.Flags().StringVar(&old, "old-flag", "", "deprecated")
cmd.Flags().MarkDeprecated("old-flag", "use --new-flag instead")
```

### Argument Validation

```go
Args: cobra.ExactArgs(1)      // Exact number
Args: cobra.RangeArgs(1, 3)   // Range
Args: cobra.MinimumNArgs(1)   // Minimum
Args: cobra.MaximumNArgs(2)   // Maximum
Args: cobra.NoArgs            // No args
Args: cobra.ArbitraryArgs     // Any number

// Custom validation
Args: func(cmd *cobra.Command, args []string) error {
    if len(args) < 1 {
        return fmt.Errorf("requires at least 1 argument")
    }
    return nil
},
```

---

## 3. Configuration (Viper)

### internal/config/config.go

```go
package config

import (
    "fmt"
    "time"

    "github.com/spf13/viper"
)

type Config struct {
    Verbose  bool   `mapstructure:"verbose"`
    LogLevel string `mapstructure:"log_level"`
    Output   string `mapstructure:"output"`
    Project  ProjectConfig `mapstructure:"project"`
    Server   ServerConfig  `mapstructure:"server"`
    API      APIConfig     `mapstructure:"api"`
}

type ProjectConfig struct {
    DefaultTemplate string `mapstructure:"default_template"`
    OutputDir       string `mapstructure:"output_dir"`
}

type ServerConfig struct {
    Host         string        `mapstructure:"host"`
    Port         int           `mapstructure:"port"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type APIConfig struct {
    BaseURL string        `mapstructure:"base_url"`
    Key     string        `mapstructure:"key"`
    Timeout time.Duration `mapstructure:"timeout"`
}

func Load() (*Config, error) {
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("unmarshaling config: %w", err)
    }
    return &cfg, nil
}

func (c *Config) Validate() error {
    if c.Server.Port < 0 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid port: %d", c.Server.Port)
    }
    return nil
}
```

### internal/config/defaults.go

```go
package config

import (
    "time"
    "github.com/spf13/viper"
)

func SetDefaults() {
    viper.SetDefault("verbose", false)
    viper.SetDefault("log_level", "info")
    viper.SetDefault("output", "text")
    viper.SetDefault("project.default_template", "default")
    viper.SetDefault("project.output_dir", ".")
    viper.SetDefault("server.host", "localhost")
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("server.read_timeout", 30*time.Second)
    viper.SetDefault("server.write_timeout", 30*time.Second)
    viper.SetDefault("api.base_url", "https://api.example.com")
    viper.SetDefault("api.timeout", 10*time.Second)
}
```

### configs/config.example.yaml

```yaml
# mycli configuration

verbose: false
log_level: info  # debug, info, warn, error
output: text     # text, json, yaml

project:
  default_template: go-api
  output_dir: ./projects

server:
  host: localhost
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

api:
  base_url: https://api.example.com
  key: ${MYCLI_API_KEY}
  timeout: 10s
```

### Configuration Priority

1. Explicit `Set()` calls (highest)
2. Flags
3. Environment variables
4. Config file
5. Defaults (lowest)

---

## 4. Logging (slog)

### internal/logger/logger.go

```go
package logger

import (
    "context"
    "io"
    "log/slog"
    "os"
)

func Init(level slog.Level, production bool) {
    opts := &slog.HandlerOptions{
        Level:     level,
        AddSource: !production && level == slog.LevelDebug,
    }

    var handler slog.Handler
    if production {
        handler = slog.NewJSONHandler(os.Stdout, opts)
    } else {
        handler = slog.NewTextHandler(os.Stderr, opts)
    }

    slog.SetDefault(slog.New(handler))
}

func InitWithWriter(w io.Writer, level slog.Level) *slog.Logger {
    opts := &slog.HandlerOptions{Level: level}
    return slog.New(slog.NewTextHandler(w, opts))
}

func FromContext(ctx context.Context) *slog.Logger {
    if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
        return l
    }
    return slog.Default()
}

func WithContext(ctx context.Context, l *slog.Logger) context.Context {
    return context.WithValue(ctx, loggerKey{}, l)
}

type loggerKey struct{}

func With(args ...any) *slog.Logger {
    return slog.Default().With(args...)
}
```

### Usage

```go
slog.Debug("processing", "id", itemID)
slog.Info("completed", "count", count)
slog.Warn("deprecated", "feature", name)
slog.Error("failed", "error", err)

// Child logger
logger := slog.With("component", "handler")
logger.Info("starting")

// Groups
slog.Info("request",
    slog.Group("user", slog.String("id", userID)),
    slog.Group("request", slog.String("method", method)),
)
```

---

## 5. Error Handling

### internal/model/errors.go

```go
package model

import (
    "errors"
    "fmt"
)

var (
    ErrNotFound      = errors.New("not found")
    ErrAlreadyExists = errors.New("already exists")
    ErrInvalidInput  = errors.New("invalid input")
    ErrUnauthorized  = errors.New("unauthorized")
)

type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

func NewValidationError(field, message string) *ValidationError {
    return &ValidationError{Field: field, Message: message}
}

type OperationError struct {
    Op     string
    Entity string
    ID     string
    Err    error
}

func (e *OperationError) Error() string {
    if e.ID != "" {
        return fmt.Sprintf("%s %s %s: %v", e.Op, e.Entity, e.ID, e.Err)
    }
    return fmt.Sprintf("%s %s: %v", e.Op, e.Entity, e.Err)
}

func (e *OperationError) Unwrap() error {
    return e.Err
}

func NewOperationError(op, entity, id string, err error) *OperationError {
    return &OperationError{Op: op, Entity: entity, ID: id, Err: err}
}
```

### Error Patterns

```go
// Wrap with context
return fmt.Errorf("reading config %s: %w", path, err)

// Check specific error
if errors.Is(err, model.ErrNotFound) {
    // Handle not found
}

// Check error type
var ve *model.ValidationError
if errors.As(err, &ve) {
    fmt.Printf("Invalid %s: %s\n", ve.Field, ve.Message)
}
```

---

## 6. Testing Patterns

### Table-Driven Tests

```go
func TestProcess(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "hello", "HELLO", false},
        {"empty input", "", "", true},
    }

    for _, tt := range tests {
        tt := tt
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

### Mocks

```go
type MockRepository struct {
    GetByIDFunc func(ctx context.Context, id string) (*model.Project, error)
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*model.Project, error) {
    if m.GetByIDFunc != nil {
        return m.GetByIDFunc(ctx, id)
    }
    return nil, nil
}
```

### Test Helpers

```go
func AssertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func AssertEqual[T comparable](t *testing.T, got, want T) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

### Running Tests

```bash
go test ./...                    # All tests
go test -v ./...                 # Verbose
go test -race ./...              # Race detection
go test -cover ./...             # Coverage
go test -tags=integration ./...  # Integration
```

---

## 7. Static Analysis Tools

### .golangci.yml

```yaml
version: "2"

run:
  timeout: 5m
  modules-download-mode: readonly

linters:
  default: standard
  enable:
    - errcheck
    - govet
    - staticcheck
    - gosimple
    - unused
    - ineffassign
    - gofmt
    - goimports
    - misspell
    - cyclop
    - gocognit
    - nestif
    - goconst
    - gocritic
    - revive
    - unconvert
    - unparam
    - errorlint
    - wrapcheck
    - nilerr
    - gosec
    - bodyclose
    - sqlclosecheck
    - tparallel
    - prealloc

  settings:
    cyclop:
      max-complexity: 15
    gocognit:
      min-complexity: 20
    nestif:
      min-complexity: 5

  exclusions:
    rules:
      - path: _test\.go
        linters:
          - dupl
          - funlen
          - gocognit
          - gosec
          - wrapcheck

issues:
  max-issues-per-linter: 50
  max-same-issues: 10
```

### Commands

```bash
golangci-lint run ./...          # Run all
golangci-lint run --fix ./...    # Auto-fix
govulncheck ./...                # Vulnerabilities
gosec ./...                      # Security
```

---

## 8. Makefile

```makefile
BINARY_NAME := mycli
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

.PHONY: release-dry
release-dry: ## Dry run release
	goreleaser release --snapshot --clean

.PHONY: release
release: ## Create release
	goreleaser release --clean
```

---

## 9. Developer Experience Files

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
/build/bin/
*.exe
*.test
coverage/
coverage.out
*.prof
.idea/
.vscode/
*.swp
.DS_Store
.env
.mycli.yaml
dist/
```

### lefthook.yml

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
      run: golangci-lint run --fix --new-from-rev=HEAD~1
      stage_fixed: true

pre-push:
  commands:
    test:
      run: go test -race -short ./...
    vet:
      run: go vet ./...
```

### pkg/version/version.go

```go
package version

import (
    "fmt"
    "runtime"
)

var (
    Version   = "dev"
    Commit    = "unknown"
    Date      = "unknown"
    Branch    = "unknown"
    GoVersion = runtime.Version()
)

type Info struct {
    Version   string `json:"version"`
    Commit    string `json:"commit"`
    Date      string `json:"date"`
    Branch    string `json:"branch"`
    GoVersion string `json:"go_version"`
    Platform  string `json:"platform"`
}

func Get() Info {
    return Info{
        Version:   Version,
        Commit:    Commit,
        Date:      Date,
        Branch:    Branch,
        GoVersion: GoVersion,
        Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
    }
}

func (i Info) String() string {
    return fmt.Sprintf("%s (commit: %s, built: %s, %s)",
        i.Version, i.Commit, i.Date, i.GoVersion)
}
```

---

## 10. CI/CD (GitHub Actions)

### .github/workflows/ci.yml

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

env:
  GO_VERSION: '1.22'

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - uses: golangci/golangci-lint-action@v6

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - run: go test -race -coverprofile=coverage.out ./...
      - uses: codecov/codecov-action@v4

  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...

  build:
    runs-on: ubuntu-latest
    needs: [lint, test]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - run: make build
      - uses: actions/upload-artifact@v4
        with:
          name: mycli
          path: build/bin/mycli
```

### .github/workflows/release.yml

```yaml
name: Release

on:
  push:
    tags: ['v*']

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - uses: goreleaser/goreleaser-action@v5
        with:
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## 11. Docker

### build/package/Dockerfile

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev COMMIT=unknown DATE=unknown
RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w -X pkg/version.Version=${VERSION}" \
    -o /mycli .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D appuser
USER appuser
COPY --from=builder /mycli /usr/local/bin/
ENTRYPOINT ["mycli"]
```

---

## 12. Releasing (GoReleaser)

### .goreleaser.yml

```yaml
version: 2
project_name: mycli

builds:
  - id: mycli
    main: .
    binary: mycli
    env: [CGO_ENABLED=0]
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    ldflags:
      - -s -w
      - -X pkg/version.Version={{.Version}}
      - -X pkg/version.Commit={{.Commit}}

archives:
  - formats: [tar.gz]
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: 'checksums.txt'

changelog:
  sort: asc
  filters:
    exclude: ['^docs:', '^test:', '^chore:']
```

---

## 13. Shell Completion

### cmd/completion.go

```go
package cmd

import (
    "os"
    "github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
    Use:   "completion [bash|zsh|fish|powershell]",
    Short: "Generate shell completion scripts",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        switch args[0] {
        case "bash":
            return rootCmd.GenBashCompletion(os.Stdout)
        case "zsh":
            return rootCmd.GenZshCompletion(os.Stdout)
        case "fish":
            return rootCmd.GenFishCompletion(os.Stdout, true)
        case "powershell":
            return rootCmd.GenPowerShellCompletion(os.Stdout)
        }
        return nil
    },
}
```

---

## 14. Profiling & Benchmarking

### Benchmarks

```go
func BenchmarkProcess(b *testing.B) {
    h := NewHandler()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = h.Process(context.Background(), "input")
    }
}
```

### Commands

```bash
go test -bench=. -benchmem ./...
go test -cpuprofile=cpu.prof -bench=. ./...
go tool pprof cpu.prof
```

---

## 15. Implementation Checklist

### Phase 1: Foundation
- [ ] `go mod init`
- [ ] `main.go`, `cmd/root.go`, `cmd/version.go`
- [ ] `.gitignore`, `.editorconfig`

### Phase 2: Core
- [ ] `internal/logger/`
- [ ] `internal/config/`
- [ ] `pkg/version/`

### Phase 3: Quality
- [ ] `.golangci.yml`
- [ ] `Makefile`
- [ ] `lefthook.yml`
- [ ] `make audit` passes

### Phase 4: Features
- [ ] Commands, handlers, services
- [ ] Tests for all code

### Phase 5: CI/CD
- [ ] GitHub Actions
- [ ] GoReleaser
- [ ] Docker

---

## 16. Code Templates

### Handler

```go
type FeatureInput struct{ Name string }
type FeatureOutput struct{ Result string }

type featureHandler struct{}

func NewFeatureHandler() *featureHandler { return &featureHandler{} }

func (h *featureHandler) Execute(ctx context.Context, in FeatureInput) (*FeatureOutput, error) {
    if in.Name == "" {
        return nil, fmt.Errorf("name required")
    }
    return &FeatureOutput{Result: in.Name}, nil
}
```

### Test

```go
func TestFeatureHandler(t *testing.T) {
    tests := []struct {
        name    string
        input   FeatureInput
        wantErr bool
    }{
        {"valid", FeatureInput{Name: "test"}, false},
        {"empty", FeatureInput{Name: ""}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            h := NewFeatureHandler()
            _, err := h.Execute(context.Background(), tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

---

## Quick Reference

```bash
# Setup
go mod init github.com/yourorg/mycli && make setup

# Development
make build    # Build
make test     # Test
make lint     # Lint
make audit    # All checks
make coverage # Coverage report

# Release
git tag v1.0.0 && git push --tags
```

---

## References

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Cobra](https://cobra.dev/)
- [Viper](https://github.com/spf13/viper)
- [golangci-lint](https://golangci-lint.run/)
- [GoReleaser](https://goreleaser.com/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
