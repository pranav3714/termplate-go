# Quick Reference

> **For**: Fast lookups of common patterns, commands, and file locations.
> **See also**: `PROJECT_CONTEXT.md` for architecture, `CONVENTIONS.md` for standards.

## File Locations Cheat Sheet

| Need to... | File | Action |
|------------|------|--------|
| Add a command | `cmd/mycommand/mycommand.go` | CREATE |
| Register command | `cmd/root.go` | EDIT (add import + AddCommand) |
| Add handler | `internal/handler/myhandler.go` | CREATE |
| Add service | `internal/service/myservice/service.go` | CREATE |
| Add config settings | `internal/config/config.go` | EDIT (add struct fields) |
| Set config defaults | `internal/config/defaults.go` | EDIT (add viper.SetDefault) |
| Add domain error | `internal/model/errors.go` | EDIT |
| Use output formatting | `internal/output/formatter.go` | USE |
| Add test | `*_test.go` (next to code) | CREATE |

## Command Snippets

### Create New Command

```go
// cmd/mycommand/mycommand.go
package mycommand

import (
    "context"
    "fmt"
    "github.com/spf13/cobra"
    "github.com/blacksilver/termplate/internal/handler"
)

var myFlag string

var Cmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Brief description",
    RunE: func(cmd *cobra.Command, _ []string) error {
        h := handler.NewMyHandler()
        return h.Execute(cmd.Context(), handler.MyInput{Field: myFlag})
    },
}

func init() {
    Cmd.Flags().StringVarP(&myFlag, "myflag", "m", "", "description")
}
```

### Register Command

```go
// cmd/root.go
import "github.com/blacksilver/termplate/cmd/mycommand"

func init() {
    // ... existing code ...
    rootCmd.AddCommand(mycommand.Cmd)
}
```

### Create Handler

```go
// internal/handler/myhandler.go
package handler

import (
    "context"
    "fmt"
    "github.com/blacksilver/termplate/internal/model"
    "github.com/blacksilver/termplate/internal/service/myservice"
)

type MyInput struct {
    Field string
}

type MyOutput struct {
    Result string
}

type MyHandler struct {
    service *myservice.Service
}

func NewMyHandler() *MyHandler {
    return &MyHandler{service: myservice.NewService()}
}

func (h *MyHandler) Execute(ctx context.Context, in MyInput) (*MyOutput, error) {
    if in.Field == "" {
        return nil, model.NewValidationError("field", "required")
    }

    result, err := h.service.DoWork(ctx, in.Field)
    if err != nil {
        return nil, fmt.Errorf("doing work: %w", err)
    }

    return &MyOutput{Result: result}, nil
}
```

### Create Service

```go
// internal/service/myservice/service.go
package myservice

import (
    "context"
    "log/slog"
)

type Service struct{}

func NewService() *Service {
    return &Service{}
}

func (s *Service) DoWork(ctx context.Context, input string) (string, error) {
    slog.InfoContext(ctx, "doing work", "input", input)
    return "result", nil
}
```

## Configuration Snippets

### Add Config Fields

```go
// internal/config/config.go
type Config struct {
    // ... existing fields ...
    MyFeature MyFeatureConfig `mapstructure:"myfeature"`
}

type MyFeatureConfig struct {
    Setting1 string `mapstructure:"setting1"`
    Setting2 int    `mapstructure:"setting2"`
}
```

### Add Defaults

```go
// internal/config/defaults.go
func SetDefaults() {
    // ... existing defaults ...
    viper.SetDefault("myfeature.setting1", "default_value")
    viper.SetDefault("myfeature.setting2", 42)
}
```

### Use Configuration

```go
// Load full config
cfg, err := config.Load()
if err != nil {
    return fmt.Errorf("loading config: %w", err)
}
value := cfg.MyFeature.Setting1

// Or use viper directly
value := viper.GetString("myfeature.setting1")
```

## Output Formatting Snippets

### Format as JSON/YAML/Table/CSV

```go
import (
    "github.com/blacksilver/termplate/internal/config"
    "github.com/blacksilver/termplate/internal/output"
)

cfg, _ := config.Load()
formatter := output.NewFormatter(cfg.Output)

// Single map
data := map[string]string{"key": "value"}
formatter.Print(data)

// Slice of maps (table)
data := []map[string]string{
    {"name": "Alice", "age": "30"},
    {"name": "Bob", "age": "25"},
}
formatter.Print(data)

// 2D array
data := [][]string{
    {"Header1", "Header2"},
    {"Value1", "Value2"},
}
formatter.Print(data)
```

## Error Handling Snippets

### Wrap Errors

```go
if err != nil {
    return fmt.Errorf("context: %w", err)
}
```

### Domain Errors

```go
// Validation error
return model.NewValidationError("email", "invalid format")

// Operation error
return model.NewOperationError("create", "user", userID, err)

// Check specific error
if errors.Is(err, model.ErrNotFound) {
    // handle not found
}

// Check error type
var ve *model.ValidationError
if errors.As(err, &ve) {
    fmt.Printf("Field %s: %s\n", ve.Field, ve.Message)
}
```

## Logging Snippets

### Structured Logging

```go
import "log/slog"

// Debug
slog.Debug("debug message", "key", value)

// Info
slog.Info("info message", "user_id", userID, "action", action)

// Warn
slog.Warn("warning", "issue", description)

// Error
slog.Error("error occurred", "error", err, "context", ctx)

// With context (includes trace IDs)
slog.InfoContext(ctx, "message", "key", value)
```

## Testing Snippets

### Table-Driven Test

```go
func TestMyFunction(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid", "input", "output", false},
        {"invalid", "", "", true},
    }

    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            got, err := MyFunction(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Make Commands

```bash
make help         # Show all commands
make build        # Build binary
make test         # Run tests
make coverage     # Coverage report
make fmt          # Format code
make lint         # Run linters
make lint-fix     # Auto-fix linting
make vet          # Run go vet
make vuln         # Check vulnerabilities
make tidy         # Tidy dependencies
make audit        # All quality checks
make ci           # Full CI pipeline
make clean        # Remove artifacts
```

## Git Workflow

```bash
# Make changes
git add .
git commit -m "feat: Add new feature"  # Auto-formats, lints
git push                                # Auto-runs tests

# Quality checks before commit
make audit
```

## Environment Variables

```bash
# General
export TERMPLATE_VERBOSE=true
export TERMPLATE_LOG_LEVEL=debug

# Output
export TERMPLATE_OUTPUT_FORMAT=json
export TERMPLATE_OUTPUT_COLOR=false

# API
export TERMPLATE_API_KEY=xxx
export TERMPLATE_API_BASE_URL=https://api.example.com

# Files
export TERMPLATE_FILES_INPUT_DIR=/path/to/input
export TERMPLATE_FILES_OUTPUT_DIR=/path/to/output

# Database
export TERMPLATE_DB_USER=user
export TERMPLATE_DB_PASSWORD=pass
export TERMPLATE_DB_HOST=localhost
export TERMPLATE_DB_PORT=5432
```

## Configuration Files

```bash
# Example config
configs/config.example.yaml

# User config locations (in priority order)
~/.termplate.yaml   # Home directory
./.termplate.yaml   # Current directory
--config /path/to/config.yaml # Via flag
```

## Common Patterns

### Context Passing

```go
// Always accept context as first parameter
func DoWork(ctx context.Context, input string) error

// Pass context through call chain
h.service.DoWork(ctx, input)
```

### Dependency Injection

```go
// Service accepts dependencies in constructor
type Service struct {
    repo repository.Interface
    client *http.Client
}

func NewService(repo repository.Interface, client *http.Client) *Service {
    return &Service{repo: repo, client: client}
}
```

### Interface Usage

```go
// Define interface near usage
type Repository interface {
    Get(ctx context.Context, id string) (*Entity, error)
}

// Accept interfaces
func NewService(repo Repository) *Service
```

## Flag Patterns

```go
// String flag
cmd.Flags().StringVarP(&var, "name", "n", "default", "description")

// Bool flag
cmd.Flags().BoolVarP(&var, "verbose", "v", false, "description")

// Int flag
cmd.Flags().IntVarP(&var, "count", "c", 10, "description")

// String slice flag
cmd.Flags().StringSliceVarP(&var, "tag", "t", []string{}, "description")

// Duration flag
cmd.Flags().DurationVarP(&var, "timeout", "T", 30*time.Second, "description")

// Required flag
cmd.MarkFlagRequired("name")

// Mutually exclusive
cmd.MarkFlagsMutuallyExclusive("flag1", "flag2")
```

## Documentation Locations

| Topic | Document |
|-------|----------|
| Project overview | `PROJECT_CONTEXT.md` |
| Coding standards | `CONVENTIONS.md` |
| Quick reference | `QUICK_REFERENCE.md` (this file) |
| Next steps | `docs/NEXT_STEPS.md` |
| Add features | `docs/GETTING_STARTED.md` |
| Configuration | `docs/CONFIGURATION_GUIDE.md` |
| CLI patterns | `docs/GO_CLI_COMPREHENSIVE_REFERENCE.md` |
| Full overview | `docs/PROJECT_SUMMARY.md` |

## Architecture Diagram

```
User Input
    ↓
cmd/                     [CLI layer - thin wiring]
    ↓
internal/handler/        [Validation + I/O formatting]
    ↓
internal/service/        [Business logic - testable]
    ↓
internal/repository/     [Data access - abstracted]
    ↓
internal/model/          [Domain entities]
```

## Example Lookups

### "How do I add a command?"
1. Create `cmd/mycommand/mycommand.go` (see snippet above)
2. Create `internal/handler/myhandler.go`
3. Create `internal/service/myservice/service.go`
4. Register in `cmd/root.go`

### "How do I add configuration?"
1. Add struct to `internal/config/config.go`
2. Add defaults in `internal/config/defaults.go`
3. Document in `configs/config.example.yaml`
4. Use with `viper.Get*()` or `config.Load()`

### "How do I format output?"
```go
cfg, _ := config.Load()
formatter := output.NewFormatter(cfg.Output)
formatter.Print(myData)
```

### "How do I log?"
```go
slog.Info("message", "key", value)
slog.InfoContext(ctx, "message", "key", value)
```

### "How do I handle errors?"
```go
if err != nil {
    return fmt.Errorf("context: %w", err)
}
```

### "Where are examples?"
- Working command: `cmd/example/greet.go`
- Working handler: `internal/handler/greet.go`
- Working service: `internal/service/example/service.go`

## Need More Details?

- **Architecture**: `PROJECT_CONTEXT.md`
- **Standards**: `CONVENTIONS.md`
- **Full guides**: `docs/` directory
- **Main README**: `README.md`
