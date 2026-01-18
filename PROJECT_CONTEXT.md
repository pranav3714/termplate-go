# Project Context for AI Analysis

> **For AI Models**: This file provides complete project context for analysis and development tasks.

## Quick Project Overview

- **Project Name**: Termplate Go
- **Type**: Production-ready Go CLI tool template
- **Module**: `github.com/blacksilver/termplate-go`
- **Binary**: `termplate`
- **Purpose**: Multi-purpose development CLI with API client, file processing, and database capabilities

## Project State

**Current Status**: Fully configured foundation ready for feature development

- ✅ Base CLI structure complete (Cobra + Viper)
- ✅ Configuration system implemented (API, Files, Database, Output)
- ✅ Output formatting (JSON, YAML, Table, CSV)
- ✅ Git repository initialized with hooks
- ✅ Quality tools configured (golangci-lint, lefthook)
- ✅ Documentation organized in `docs/`
- ✅ Example command demonstrating full architecture
- 🔨 Ready for feature implementation

## Architecture Overview

### Layer Structure (Clean Architecture)

```
main.go                     # Entry point (minimal)
    ↓
cmd/                        # CLI layer (Cobra commands)
    ↓
internal/handler/           # Validation & I/O formatting
    ↓
internal/service/           # Business logic (testable)
    ↓
internal/repository/        # Data access (abstracted)
    ↓
internal/model/             # Domain entities (no dependencies)
```

**Key Principle**: Dependencies only flow downward, no circular imports.

### Directory Structure

```
termplate/
├── cmd/                    # CLI commands (thin wiring only)
│   ├── root.go            # Root command + global flags
│   ├── version.go         # Version command
│   ├── completion.go      # Shell completion
│   └── example/           # Example command group
│
├── internal/              # Private application code
│   ├── config/           # Configuration (Viper)
│   │   ├── config.go     # Config structs
│   │   └── defaults.go   # Default values
│   ├── logger/           # Logging (slog)
│   ├── model/            # Domain models + errors
│   ├── handler/          # Command handlers
│   ├── service/          # Business logic
│   ├── output/           # Output formatting
│   ├── repository/       # Data access (empty, ready to use)
│   ├── validator/        # Input validation (empty, ready to use)
│   └── testutil/         # Test helpers (empty, ready to use)
│
├── pkg/                   # Public packages (importable)
│   └── version/          # Version information
│
├── configs/               # Configuration templates
│   └── config.example.yaml
│
├── docs/                  # All documentation
│   ├── README.md         # Documentation index
│   ├── NEXT_STEPS.md     # What to do next
│   ├── GETTING_STARTED.md # How to add features
│   ├── CONFIGURATION_GUIDE.md # Config reference
│   └── ...more
│
├── test/                  # Integration & E2E tests
│   ├── integration/
│   ├── e2e/
│   └── testdata/
│
└── build/                 # Build artifacts
    ├── bin/              # Compiled binary
    └── package/          # Docker files
```

## Key Files for Common Tasks

### Adding a New Command

**Primary Files**:
- `cmd/mycommand/mycommand.go` - Command definition (CREATE NEW)
- `cmd/root.go` - Register command (MODIFY)
- `internal/handler/mycommand.go` - Handler logic (CREATE NEW)
- `internal/service/mycommand/service.go` - Business logic (CREATE NEW)

**Reference**: `docs/GETTING_STARTED.md` (section: "Adding Your First Command")

### Configuration

**Primary Files**:
- `internal/config/config.go` - Config structures
- `internal/config/defaults.go` - Default values
- `configs/config.example.yaml` - Example configuration
- `cmd/root.go` - Config initialization

**Reference**: `docs/CONFIGURATION_GUIDE.md`

### Error Handling

**Primary Files**:
- `internal/model/errors.go` - Domain errors

**Pattern**: Always wrap errors with context using `fmt.Errorf("context: %w", err)`

### Logging

**Primary Files**:
- `internal/logger/logger.go` - Logger setup

**Pattern**: Use structured logging: `slog.Info("message", "key", value)`

### Output Formatting

**Primary Files**:
- `internal/output/formatter.go` - Multi-format output

**Usage**: `formatter.Print(data)` automatically formats based on config

## Configuration Structure

### Available Configuration Sections

1. **Output** (`output.*`)
   - Format: text, json, yaml, table, csv
   - Pretty printing, colors, table styles

2. **API** (`api.*`)
   - Base URL, keys, tokens
   - Retry logic, timeouts
   - Custom headers, rate limiting

3. **Files** (`files.*`)
   - Input/output directories
   - Patterns, size limits
   - Backup, permissions

4. **Database** (`database.*`)
   - Driver: postgres, mysql, sqlite
   - Connection pooling
   - SSL mode, timeouts

5. **Server** (`server.*`)
   - Host, port
   - Timeouts, TLS

### Environment Variables

Prefix: `TERMPLATE_`

Examples:
- `TERMPLATE_API_KEY`
- `TERMPLATE_OUTPUT_FORMAT`
- `TERMPLATE_DB_PASSWORD`

## Development Workflow

### Before Coding

1. Read relevant docs in `docs/`
2. Check existing patterns in `cmd/example/`
3. Review configuration in `configs/config.example.yaml`

### When Adding Features

1. **Explore**: Check existing code structure
2. **Plan**: Follow clean architecture layers
3. **Code**: Implement cmd → handler → service → repository
4. **Test**: Write table-driven tests
5. **Format**: Run `make fmt`
6. **Lint**: Run `make lint`
7. **Commit**: Git hooks will auto-format and lint

### Quality Checks

```bash
make fmt        # Format code
make lint       # Run linters (20+ linters)
make test       # Run tests
make audit      # All quality checks
```

## Code Patterns

### Command Pattern (cmd/)

```go
package mycommand

import (
    "github.com/spf13/cobra"
    "github.com/blacksilver/termplate/internal/handler"
)

var Cmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Brief description",
    RunE: func(cmd *cobra.Command, _ []string) error {
        h := handler.NewMyHandler()
        return h.Execute(cmd.Context(), handler.MyInput{})
    },
}
```

### Handler Pattern (internal/handler/)

```go
package handler

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
    // Validate input
    if in.Field == "" {
        return nil, model.NewValidationError("field", "field is required")
    }

    // Call service
    result, err := h.service.DoSomething(ctx, in.Field)
    if err != nil {
        return nil, fmt.Errorf("doing something: %w", err)
    }

    return &MyOutput{Result: result}, nil
}
```

### Service Pattern (internal/service/)

```go
package myservice

type Service struct {
    // Dependencies (repositories, clients, etc.)
}

func NewService() *Service {
    return &Service{}
}

func (s *Service) DoSomething(ctx context.Context, input string) (string, error) {
    slog.InfoContext(ctx, "doing something", "input", input)

    // Business logic here
    return "result", nil
}
```

### Test Pattern

```go
func TestService_DoSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid", "test", "result", false},
        {"empty", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := NewService()
            got, err := s.DoSomething(context.Background(), tt.input)

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

## Important Conventions

### Rules

1. **cmd/** - Only CLI wiring, no business logic
2. **internal/handler/** - Thin layer: validation + I/O formatting
3. **internal/service/** - All business logic (framework-independent)
4. **Dependencies** - Only flow downward, never circular
5. **Errors** - Always wrap with context: `fmt.Errorf("context: %w", err)`
6. **Logging** - Structured logging: `slog.Info("msg", "key", value)`
7. **Tests** - Table-driven tests for all logic
8. **Exports** - Only export types/functions that need to be public

### Naming Conventions

- **Commands**: lowercase, hyphenated (e.g., `my-command`)
- **Packages**: lowercase, single word (e.g., `mypackage`)
- **Files**: lowercase, underscored (e.g., `my_handler.go`)
- **Types**: PascalCase (e.g., `MyHandler`)
- **Functions**: camelCase (e.g., `doSomething`)
- **Constants**: PascalCase or UPPER_SNAKE (e.g., `MaxRetries` or `MAX_RETRIES`)

## Dependencies

### Core Framework

- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration
- `gopkg.in/yaml.v3` - YAML support
- Standard library `log/slog` - Logging

### Development Tools

- `golangci-lint` - Linting (20+ linters enabled)
- `lefthook` - Git hooks
- `govulncheck` - Vulnerability scanning
- `goimports` - Import formatting

## Common Mistakes to Avoid

1. ❌ **Don't** put business logic in `cmd/`
2. ❌ **Don't** create circular imports
3. ❌ **Don't** ignore error wrapping
4. ❌ **Don't** use `panic()` for normal errors
5. ❌ **Don't** forget to validate input in handlers
6. ❌ **Don't** make everything public (export only what's needed)
7. ❌ **Don't** skip writing tests
8. ❌ **Don't** commit without running `make lint`

## Quick Command Reference

### Build & Run

```bash
make build                              # Build binary
./build/bin/termplate --help  # Run CLI
./build/bin/termplate version # Check version
```

### Development

```bash
make fmt        # Format code
make lint       # Run linters
make test       # Run tests
make coverage   # Generate coverage report
make audit      # Run all quality checks
```

### Git Workflow

```bash
git add .
git commit -m "message"  # Auto-formats and lints
git push                 # Auto-runs tests
```

## Documentation Map

**For specific tasks, consult these docs:**

| Task | Document | Location |
|------|----------|----------|
| Getting started | Getting Started Guide | `docs/GETTING_STARTED.md` |
| Next steps | Next Steps | `docs/NEXT_STEPS.md` |
| Configuration | Configuration Guide | `docs/CONFIGURATION_GUIDE.md` |
| CLI patterns | CLI Reference | `docs/GO_CLI_COMPREHENSIVE_REFERENCE.md` |
| Project overview | Project Summary | `docs/PROJECT_SUMMARY.md` |
| What's customized | Customization Complete | `docs/CUSTOMIZATION_COMPLETE.md` |

## AI Model Analysis Tips

When analyzing this project:

1. **Start with** `PROJECT_CONTEXT.md` (this file) for overview
2. **Check** `docs/` for specific guidance
3. **Review** `cmd/example/` for working examples
4. **Follow** the architecture: cmd → handler → service → repository
5. **Reference** `internal/config/config.go` for configuration
6. **Use** `internal/output/formatter.go` for output formatting
7. **Check** `internal/model/errors.go` for error patterns
8. **Follow** conventions in `CONVENTIONS.md`

## Example Workflow: Adding a File Processing Command

1. **Read**: `docs/GETTING_STARTED.md` (section on adding commands)
2. **Check config**: `internal/config/config.go` (FilesConfig section)
3. **Create command**: `cmd/process/process.go`
4. **Create handler**: `internal/handler/process.go`
5. **Create service**: `internal/service/process/service.go`
6. **Register**: Add to `cmd/root.go`
7. **Test**: Write tests for service
8. **Verify**: Run `make audit`

## Build Information

- **Go Version**: 1.22+
- **Module**: `github.com/blacksilver/termplate`
- **Binary**: `termplate`
- **Build Command**: `make build`
- **Output**: `./build/bin/termplate`

## Contact & Resources

- **Documentation**: `docs/` directory
- **Examples**: `cmd/example/` directory
- **Configuration**: `configs/config.example.yaml`
- **Main README**: `README.md`

---

**Last Updated**: 2026-01-18
**Status**: Foundation complete, ready for feature development
