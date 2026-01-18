# Project Conventions

> **Purpose**: Defines coding standards, patterns, and conventions for this project.
> **Audience**: Developers and AI models working on this codebase.

## Table of Contents

- [Architecture](#architecture)
- [Directory Structure](#directory-structure)
- [Naming Conventions](#naming-conventions)
- [Code Patterns](#code-patterns)
- [Error Handling](#error-handling)
- [Logging](#logging)
- [Testing](#testing)
- [Configuration](#configuration)
- [Documentation](#documentation)

## Architecture

### Clean Architecture Principles

This project follows clean architecture with strict layer separation:

```
Presentation Layer:  cmd/                (CLI interface)
Application Layer:   internal/handler/   (orchestration, validation)
Domain Layer:        internal/service/   (business logic)
Infrastructure:      internal/repository/ (data access)
Entities:           internal/model/      (domain models)
```

### Dependency Rules

1. **Unidirectional Flow**: Dependencies only flow inward (downward)
2. **No Circular Imports**: Never create circular dependencies
3. **Interface Boundaries**: Use interfaces at layer boundaries
4. **Framework Independence**: Core logic (service layer) has no framework dependencies

### Layer Responsibilities

| Layer | Responsibility | Should NOT Do |
|-------|---------------|---------------|
| `cmd/` | CLI wiring, flag parsing | Business logic, data access |
| `handler/` | Input validation, output formatting | Business logic, complex operations |
| `service/` | Business logic, orchestration | Direct I/O, framework-specific code |
| `repository/` | Data access, external APIs | Business logic, validation |
| `model/` | Domain entities, errors | Dependencies on other layers |

## Directory Structure

### Standard Locations

```
termplate/
├── cmd/                    # Commands ONLY (thin layer)
│   ├── root.go            # Root command + flags
│   ├── <feature>/         # Feature command groups
│   │   ├── <feature>.go   # Parent command
│   │   └── <action>.go    # Subcommands
│
├── internal/              # Private application code
│   ├── config/           # Configuration (Viper)
│   ├── handler/          # Command handlers (thin)
│   ├── service/          # Business logic (core)
│   ├── repository/       # Data access
│   ├── model/            # Domain models
│   ├── logger/           # Logging setup
│   ├── output/           # Output formatting
│   ├── validator/        # Input validation
│   └── testutil/         # Test helpers
│
├── pkg/                   # Public, reusable packages
├── configs/               # Configuration templates
├── docs/                  # All documentation
├── test/                  # Integration tests
└── build/                 # Build artifacts
```

### File Organization

- **One type per file**: Each struct/interface in its own file (when significant)
- **Group related functions**: Keep related functions together
- **Separate concerns**: Don't mix business logic with I/O
- **Test files**: Place `*_test.go` files next to the code they test

## Naming Conventions

### Files

```
lowercase_with_underscores.go  # Standard files
mytype.go                      # File named after main type
mytype_test.go                 # Test file
```

### Packages

```
package mypackage  # lowercase, single word, descriptive
package auth       # NOT: authentication, authService
package config     # NOT: cfg, configuration
```

### Types

```go
type MyHandler struct {}      # Exported: PascalCase
type myInternalType struct {} # Unexported: camelCase
type UserID string            # Type aliases: PascalCase
```

### Functions/Methods

```go
func NewHandler() *Handler    // Constructors: New + Type
func Execute()                # Exported: PascalCase
func doInternal()             # Unexported: camelCase
func (h *Handler) Execute()   # Methods: PascalCase or camelCase
```

### Variables

```go
var ConfigFile string         # Exported: PascalCase
var maxRetries int            # Unexported: camelCase
const MaxConnections = 100    # Constants: PascalCase
const defaultTimeout = 30     # Internal: camelCase
```

### Command Names

```
termplate process           # lowercase
termplate process-file      # hyphenated
termplate api fetch-data    # nested with hyphens
```

## Code Patterns

### Command Pattern

Every command follows this structure:

```go
package mycommand

import (
    "context"
    "github.com/spf13/cobra"
    "github.com/blacksilver/termplate/internal/handler"
)

// Flag variables
var (
    flagName   string
    flagCount  int
)

// Cmd is the command definition
var Cmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Brief description",
    Long:  `Longer description with examples`,
    Args:  cobra.NoArgs,
    RunE: func(cmd *cobra.Command, _ []string) error {
        return runMyCommand(cmd.Context())
    },
}

func init() {
    Cmd.Flags().StringVarP(&flagName, "name", "n", "", "description")
    _ = Cmd.MarkFlagRequired("name")
}

func runMyCommand(ctx context.Context) error {
    h := handler.NewMyCommandHandler()
    result, err := h.Execute(ctx, handler.MyCommandInput{
        Name: flagName,
    })
    if err != nil {
        return fmt.Errorf("executing command: %w", err)
    }

    fmt.Println(result.Message)
    return nil
}
```

### Handler Pattern

```go
package handler

// Input/Output types
type MyInput struct {
    Field string
}

type MyOutput struct {
    Result string
}

// Handler type (exported)
type MyHandler struct {
    service *myservice.Service
}

// Constructor
func NewMyHandler() *MyHandler {
    return &MyHandler{
        service: myservice.NewService(),
    }
}

// Execute method
func (h *MyHandler) Execute(ctx context.Context, in MyInput) (*MyOutput, error) {
    // 1. Validate input
    if in.Field == "" {
        return nil, model.NewValidationError("field", "field is required")
    }

    // 2. Call service
    result, err := h.service.DoWork(ctx, in.Field)
    if err != nil {
        return nil, fmt.Errorf("doing work: %w", err)
    }

    // 3. Return output
    return &MyOutput{Result: result}, nil
}
```

### Service Pattern

```go
package myservice

type Service struct {
    // Dependencies (repositories, other services)
    repo repository.Interface
}

func NewService() *Service {
    return &Service{
        // Initialize dependencies
    }
}

func (s *Service) DoWork(ctx context.Context, input string) (string, error) {
    // Business logic here
    slog.InfoContext(ctx, "doing work", "input", input)

    // Return result
    return "result", nil
}
```

### Repository Pattern

```go
package myrepository

// Interface defines contract
type Interface interface {
    Get(ctx context.Context, id string) (*model.Entity, error)
    Save(ctx context.Context, entity *model.Entity) error
}

// Implementation
type repository struct {
    // Dependencies (DB, API client, etc.)
}

func New() Interface {
    return &repository{}
}

func (r *repository) Get(ctx context.Context, id string) (*model.Entity, error) {
    // Data access logic
    return nil, nil
}
```

## Error Handling

### Rules

1. **Always wrap errors** with context using `fmt.Errorf("context: %w", err)`
2. **Never ignore errors** - always handle or return them
3. **Use domain errors** from `internal/model/errors.go` when appropriate
4. **Return early** on errors (guard clauses)
5. **Don't panic** for normal error conditions

### Error Wrapping

```go
// ✅ GOOD: Wrap with context
if err := doSomething(); err != nil {
    return fmt.Errorf("doing something: %w", err)
}

// ❌ BAD: No context
if err != nil {
    return err
}

// ❌ BAD: Loses error chain
if err != nil {
    return fmt.Errorf("doing something: %v", err) // %v not %w
}
```

### Domain Errors

```go
// Check specific errors
if errors.Is(err, model.ErrNotFound) {
    // Handle not found
}

// Check error types
var ve *model.ValidationError
if errors.As(err, &ve) {
    fmt.Printf("Invalid %s: %s\n", ve.Field, ve.Message)
}

// Create domain errors
return model.NewValidationError("email", "invalid email format")
return model.NewOperationError("create", "user", userID, err)
```

## Logging

### Rules

1. **Use structured logging** with `slog`
2. **Use InfoContext/ErrorContext** to include trace context
3. **Log at appropriate levels**: Debug < Info < Warn < Error
4. **Include relevant context** as key-value pairs
5. **Don't log sensitive data** (passwords, tokens, etc.)

### Logging Levels

```go
// Debug: Detailed information for debugging
slog.Debug("processing item", "id", itemID, "count", count)

// Info: General informational messages
slog.Info("operation started", "user", userID)

// Warn: Warning messages (degraded state)
slog.Warn("slow query detected", "duration", duration)

// Error: Error conditions
slog.Error("operation failed", "error", err, "user", userID)
```

### Context-Aware Logging

```go
// Use Context variants to include trace IDs
slog.InfoContext(ctx, "processing request", "path", path)
slog.ErrorContext(ctx, "request failed", "error", err)
```

### Structured Data

```go
// ✅ GOOD: Structured with key-value pairs
slog.Info("user created",
    "user_id", userID,
    "email", email,
    "role", role,
)

// ❌ BAD: String interpolation
slog.Info(fmt.Sprintf("user %s created with email %s", userID, email))
```

## Testing

### Test Organization

```go
func TestServiceName_MethodName(t *testing.T) {
    t.Parallel() // Run in parallel when possible

    tests := []struct {
        name    string       // Test case name
        input   InputType    // Input data
        want    OutputType   // Expected output
        wantErr bool         // Expect error?
    }{
        {
            name:    "valid input",
            input:   validInput,
            want:    expectedOutput,
            wantErr: false,
        },
        {
            name:    "invalid input",
            input:   invalidInput,
            want:    OutputType{},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        tt := tt // Capture range variable
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            // Arrange
            s := NewService()

            // Act
            got, err := s.Method(context.Background(), tt.input)

            // Assert
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Test Helpers

```go
// Use t.Helper() in helper functions
func assertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func assertEqual[T comparable](t *testing.T, got, want T) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

## Configuration

### Configuration Access

```go
// Load full configuration
cfg, err := config.Load()

// Access directly via viper
apiKey := viper.GetString("api.key")
timeout := viper.GetDuration("api.timeout")
patterns := viper.GetStringSlice("files.patterns")
```

### Configuration Keys

Use dot notation for nested configuration:

```
verbose
log_level
output.format
output.color
api.base_url
api.timeout
files.input_dir
database.driver
```

### Environment Variables

Prefix: `TERMPLATE_`

```bash
TERMPLATE_API_KEY=xxx
TERMPLATE_OUTPUT_FORMAT=json
TERMPLATE_DB_PASSWORD=xxx
```

## Documentation

### Code Comments

```go
// Package-level doc
// Package mypackage provides functionality for X.
package mypackage

// Type documentation
// MyHandler handles command execution.
type MyHandler struct {}

// Function documentation
// NewMyHandler creates a new handler instance.
func NewMyHandler() *MyHandler {}

// Method documentation
// Execute runs the command with the given input.
func (h *MyHandler) Execute(ctx context.Context, in Input) (*Output, error) {}
```

### README Files

- **Root README.md**: Project overview, quick start
- **docs/README.md**: Documentation index
- **Package README**: If package is complex

### Documentation Standards

1. **Keep docs in sync** with code
2. **Use examples** liberally
3. **Document why**, not just what
4. **Include common pitfalls**
5. **Link related docs**

## Command-Line Interface

### Flag Naming

```go
// Short flags: single letter
--name, -n
--verbose, -v
--output, -o

// Long flags: lowercase with hyphens
--max-retries
--output-dir
--api-key

// Boolean flags: positive form
--enable-cache    // NOT --disable-cache
--verbose         // NOT --quiet (use as separate flag)
```

### Command Naming

```
termplate command              # Single word
termplate parent child         # Nested
termplate multi-word-command   # Hyphenated
```

### Help Text

```go
var Cmd = &cobra.Command{
    Use:   "command [flags]",
    Short: "Brief one-line description",
    Long: `Longer description with more details.

Explains what the command does, when to use it,
and any important notes.`,
    Example: `  termplate command --flag value
  termplate command --another-flag`,
}
```

## Git Workflow

### Commit Messages

```
Format: <type>: <description>

Types:
- feat: New feature
- fix: Bug fix
- docs: Documentation changes
- refactor: Code refactoring
- test: Test additions/changes
- chore: Build process, dependencies

Examples:
feat: Add file processing command
fix: Handle empty input in greet command
docs: Update configuration guide
refactor: Extract validation to separate package
```

### Hooks

- **Pre-commit**: Auto-format (`gofmt`, `goimports`), lint
- **Pre-push**: Run tests, vet

## Performance Considerations

1. **Avoid premature optimization**: Optimize only when needed
2. **Use buffered I/O**: When reading/writing large files
3. **Context cancellation**: Respect context cancellation
4. **Connection pooling**: For database/HTTP clients
5. **Lazy initialization**: Initialize expensive resources when needed

## Security Guidelines

1. **Never commit secrets**: Use environment variables
2. **Validate all input**: Especially user-provided data
3. **Use HTTPS**: For API calls (unless explicitly configured)
4. **Sanitize file paths**: Prevent directory traversal
5. **Use crypto/rand**: Not math/rand for security

## Questions?

- Check `docs/` for detailed guides
- Review `cmd/example/` for working examples
- See `PROJECT_CONTEXT.md` for architecture overview
