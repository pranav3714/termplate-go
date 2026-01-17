# Next Steps Guide 🚀

Your `ever-so-powerful-go` CLI is fully set up, customized, and ready for development!

## What's Been Done ✅

### 1. Project Customization
- ✅ Module path: `github.com/blacksilver/ever-so-powerful`
- ✅ Binary name: `ever-so-powerful-go`
- ✅ CLI descriptions updated
- ✅ Config paths updated
- ✅ Environment variables: `EVER_SO_POWERFUL_GO_*`

### 2. Git Repository
- ✅ Git initialized on `main` branch
- ✅ Lefthook hooks installed
  - Pre-commit: Auto-format and lint
  - Pre-push: Run tests and vet
- ✅ Initial commit created with clean code

### 3. Code Quality
- ✅ All linting issues fixed
- ✅ Proper error wrapping
- ✅ Exported types correctly
- ✅ Unused parameters handled
- ✅ 20+ linters passing

### 4. Quality Checks Passed
- ✅ `go mod tidy`
- ✅ `go mod verify`
- ✅ `go vet`
- ✅ `golangci-lint`
- ✅ `go test`

## Your Project Structure

```
ever-so-powerful/
├── build/bin/
│   └── ever-so-powerful-go        ← Your compiled binary
├── cmd/
│   ├── root.go                    ← Main command (customized)
│   ├── version.go                 ← Version info
│   ├── completion.go              ← Shell completion
│   └── example/                   ← Example command (kept as reference)
│       ├── example.go
│       └── greet.go
├── internal/
│   ├── handler/                   ← Command handlers (thin layer)
│   │   └── greet.go
│   ├── service/                   ← Business logic
│   │   └── example/
│   │       └── service.go
│   ├── config/                    ← Configuration
│   ├── logger/                    ← Logging
│   └── model/                     ← Domain models & errors
├── pkg/version/                   ← Version information
├── configs/                       ← Configuration templates
└── Documentation files
```

## Quick Test Your CLI

```bash
# Build (if not already built)
make build

# Basic commands
./build/bin/ever-so-powerful-go --help
./build/bin/ever-so-powerful-go version
./build/bin/ever-so-powerful-go example greet --name "Developer"

# With flags
./build/bin/ever-so-powerful-go -v example greet --name "User" --uppercase
./build/bin/ever-so-powerful-go version --output json
```

## What You Can Do Now

### Option 1: Add Your First Real Command

Create a new command following the example pattern:

```bash
# 1. Create command directory
mkdir -p cmd/mycommand

# 2. Create handler
# Edit: internal/handler/mycommand.go

# 3. Create service
mkdir -p internal/service/mycommand
# Edit: internal/service/mycommand/service.go

# 4. Wire it up in cmd/root.go
```

See `docs/GETTING_STARTED.md` for detailed instructions.

### Option 2: Customize Configuration

Update for your needs:
- `internal/config/config.go` - Add your config structures
- `internal/config/defaults.go` - Set default values
- `configs/config.example.yaml` - Document options

### Option 3: Remove Example Command

If you want a clean slate:

```bash
# Remove example files
rm -rf cmd/example
rm internal/handler/greet.go
rm -rf internal/service/example

# Edit cmd/root.go to remove:
# - import "github.com/blacksilver/ever-so-powerful/cmd/example"
# - rootCmd.AddCommand(example.Cmd)

# Rebuild
make build
```

### Option 4: Set Up GitHub

```bash
# Create repository on GitHub first, then:
git remote add origin https://github.com/blacksilver/ever-so-powerful.git
git push -u origin main

# GitHub Actions workflows are already configured!
# - .github/workflows/ci.yml (runs on push/PR)
# - .github/workflows/release.yml (runs on tags)
```

### Option 5: Add Tests

Create your first test file:

```bash
# Create test file
# internal/service/example/service_test.go

# Follow table-driven test pattern
# See docs/GETTING_STARTED.md for examples

# Run tests
go test ./...
make coverage  # With coverage report
```

## Development Workflow

```bash
# 1. Make changes to code

# 2. Git will automatically:
#    - Format code on commit (gofmt)
#    - Run linter on commit (golangci-lint)
#    - Run tests on push

# 3. Or run manually:
make fmt        # Format code
make lint       # Run linter
make test       # Run tests
make build      # Build binary

# 4. All at once:
make audit      # Runs tidy, vet, lint, test
```

## Available Make Commands

```bash
make help         # Show all commands
make build        # Build binary
make build-all    # Build for all platforms
make test         # Run tests
make coverage     # Generate coverage report
make fmt          # Format code
make vet          # Run go vet
make lint         # Run linter
make lint-fix     # Auto-fix linting issues
make vuln         # Check vulnerabilities
make tidy         # Tidy dependencies
make audit        # Run all quality checks
make ci           # Full CI pipeline
make clean        # Remove build artifacts
```

## Common Tasks

### Add a Flag to a Command

```go
var myFlag string

func init() {
    myCmd.Flags().StringVarP(&myFlag, "myflag", "m", "default", "description")
    myCmd.MarkFlagRequired("myflag")
}
```

### Use Configuration

```go
import "github.com/spf13/viper"

value := viper.GetString("my.config.key")
```

### Add Structured Logging

```go
import "log/slog"

slog.Info("operation started", "user", username, "action", action)
slog.Error("operation failed", "error", err)
```

### Handle Errors Properly

```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("doing something: %w", err)
}

// Check specific errors
if errors.Is(err, model.ErrNotFound) {
    // Handle not found
}

// Check error types
var ve *model.ValidationError
if errors.As(err, &ve) {
    fmt.Printf("Invalid %s: %s\n", ve.Field, ve.Message)
}
```

## Example: Adding Your First Command

Let's say you want to add a `process` command:

### 1. Create Command File

`cmd/process/process.go`:
```go
package process

import (
    "context"
    "fmt"
    "github.com/spf13/cobra"
    "github.com/blacksilver/ever-so-powerful/internal/handler"
)

var fileName string

var Cmd = &cobra.Command{
    Use:   "process",
    Short: "Process a file",
    RunE: func(cmd *cobra.Command, _ []string) error {
        h := handler.NewProcessHandler()
        result, err := h.Execute(cmd.Context(), handler.ProcessInput{
            FileName: fileName,
        })
        if err != nil {
            return fmt.Errorf("processing file: %w", err)
        }
        fmt.Println(result.Message)
        return nil
    },
}

func init() {
    Cmd.Flags().StringVarP(&fileName, "file", "f", "", "file to process (required)")
    _ = Cmd.MarkFlagRequired("file")
}
```

### 2. Create Handler

`internal/handler/process.go`:
```go
package handler

import (
    "context"
    "fmt"
    "github.com/blacksilver/ever-so-powerful/internal/model"
    "github.com/blacksilver/ever-so-powerful/internal/service/process"
)

type ProcessInput struct {
    FileName string
}

type ProcessOutput struct {
    Message string
}

// ProcessHandler handles file processing
type ProcessHandler struct {
    service *process.Service
}

// NewProcessHandler creates a new process handler
func NewProcessHandler() *ProcessHandler {
    return &ProcessHandler{
        service: process.NewService(),
    }
}

// Execute processes a file
func (h *ProcessHandler) Execute(ctx context.Context, in ProcessInput) (*ProcessOutput, error) {
    if in.FileName == "" {
        return nil, model.NewValidationError("file", "file name is required")
    }

    result, err := h.service.ProcessFile(ctx, in.FileName)
    if err != nil {
        return nil, fmt.Errorf("processing file: %w", err)
    }

    return &ProcessOutput{Message: result}, nil
}
```

### 3. Create Service

`internal/service/process/service.go`:
```go
package process

import (
    "context"
    "fmt"
    "log/slog"
)

// Service handles processing logic
type Service struct {
    // Add dependencies here
}

// NewService creates a new process service
func NewService() *Service {
    return &Service{}
}

// ProcessFile processes a file
func (s *Service) ProcessFile(ctx context.Context, fileName string) (string, error) {
    slog.InfoContext(ctx, "processing file", "file", fileName)

    // Your business logic here
    return fmt.Sprintf("Processed file: %s", fileName), nil
}
```

### 4. Register Command

In `cmd/root.go`:
```go
import (
    // ... existing imports ...
    "github.com/blacksilver/ever-so-powerful/cmd/process"
)

func init() {
    // ... existing code ...
    rootCmd.AddCommand(process.Cmd)
}
```

### 5. Build and Test

```bash
make build
./build/bin/ever-so-powerful-go process --file myfile.txt
```

## Configuration Example

Create `~/.ever-so-powerful-go.yaml`:

```yaml
verbose: true
log_level: debug

# Your custom config
myapp:
  setting1: value1
  setting2: value2
```

Update `internal/config/config.go`:

```go
type Config struct {
    Verbose  bool      `mapstructure:"verbose"`
    LogLevel string    `mapstructure:"log_level"`
    MyApp    MyAppConfig `mapstructure:"myapp"`
}

type MyAppConfig struct {
    Setting1 string `mapstructure:"setting1"`
    Setting2 string `mapstructure:"setting2"`
}
```

Update `internal/config/defaults.go`:

```go
func SetDefaults() {
    // ... existing defaults ...
    viper.SetDefault("myapp.setting1", "default1")
    viper.SetDefault("myapp.setting2", "default2")
}
```

## Releasing

### Create a Release

```bash
# 1. Tag a version
git tag -a v0.1.0 -m "Initial release"

# 2. Push the tag
git push origin v0.1.0

# 3. If you have GitHub Actions set up, it will automatically:
#    - Build for Linux, macOS, Windows (amd64, arm64)
#    - Create a GitHub release
#    - Upload binaries
```

### Manual Release (without GitHub Actions)

```bash
# Install goreleaser
go install github.com/goreleaser/goreleaser@latest

# Dry run
make release-dry

# Create release
make release
```

## Shell Completion

```bash
# Bash
./build/bin/ever-so-powerful-go completion bash > /etc/bash_completion.d/ever-so-powerful-go

# Zsh
./build/bin/ever-so-powerful-go completion zsh > "${fpath[1]}/_ever-so-powerful-go"

# Fish
./build/bin/ever-so-powerful-go completion fish > ~/.config/fish/completions/ever-so-powerful-go.fish
```

## Docker

```bash
# Build image
docker build -f build/package/Dockerfile -t ever-so-powerful-go:latest .

# Run
docker run --rm ever-so-powerful-go:latest version
docker run --rm ever-so-powerful-go:latest process --file test.txt
```

## Documentation

- **docs/GETTING_STARTED.md** - Detailed guide for adding features
- **docs/CUSTOMIZATION_COMPLETE.md** - What was customized
- **docs/PROJECT_SUMMARY.md** - Complete project overview
- **README.md** - General documentation
- **GO_CLI_COMPREHENSIVE_REFERENCE.md** - Full reference guide

## Tips

1. **Keep `cmd/` thin** - Only CLI wiring, delegate to handlers
2. **Business logic in `service/`** - Testable, framework-independent
3. **Use structured logging** - `slog.Info("msg", "key", value)`
4. **Wrap errors** - Always add context with `fmt.Errorf`
5. **Write tests** - Follow table-driven test patterns
6. **Run `make audit`** - Before committing major changes

## Need Help?

1. Check the documentation files
2. Look at the example command in `cmd/example/`
3. Follow existing patterns in the codebase
4. Review `GO_CLI_COMPREHENSIVE_REFERENCE.md`

## You're All Set! 🎉

Your CLI is production-ready with:
- ✅ Clean architecture
- ✅ Best practices
- ✅ Quality checks
- ✅ Git hooks
- ✅ CI/CD ready
- ✅ Documentation

Start building your features now!

```bash
./build/bin/ever-so-powerful-go --help
```

Happy coding! 🚀
