# Getting Started with mycli

This guide will help you get started with the mycli template and customize it for your needs.

## Quick Start

### 1. Test the Template

```bash
# Build the binary
make build

# Run the CLI
./build/bin/mycli --help

# Test the example command
./build/bin/mycli example greet --name "World"

# Test with verbose logging
./build/bin/mycli -v example greet --name "World"

# Check version
./build/bin/mycli version
```

### 2. Customize for Your Project

Replace `github.com/yourorg/mycli` with your actual module path:

```bash
# Update go.mod
sed -i 's|github.com/yourorg/mycli|github.com/yourusername/yourproject|g' go.mod

# Update all Go files
find . -name "*.go" -type f -exec sed -i 's|github.com/yourorg/mycli|github.com/yourusername/yourproject|g' {} +

# Update Makefile
sed -i 's|mycli|yourproject|g' Makefile

# Update .goreleaser.yml
sed -i 's|github.com/yourorg/mycli|github.com/yourusername/yourproject|g' .goreleaser.yml
sed -i 's|mycli|yourproject|g' .goreleaser.yml

# Update Dockerfile
sed -i 's|github.com/yourorg/mycli|github.com/yourusername/yourproject|g' build/package/Dockerfile

# Tidy dependencies
go mod tidy
```

### 3. Rename the Binary

Update these files to change the binary name from `mycli` to your desired name:

- `Makefile` - Change `BINARY_NAME := mycli`
- `.goreleaser.yml` - Change `project_name: mycli` and `binary: mycli`
- `README.md` - Replace all instances of `mycli`
- `cmd/root.go` - Update the `Use` field and descriptions
- `cmd/version.go` - Update the output format
- `configs/config.example.yaml` - Update the comment

### 4. Update Descriptions

Edit `cmd/root.go` to customize your CLI description:

```go
var rootCmd = &cobra.Command{
    Use:   "yourproject",
    Short: "Your brief description here",
    Long: `Your longer description here.

Examples:
  yourproject --help
  yourproject version
  yourproject do-something --flag value`,
}
```

## Project Structure Overview

```
mycli/
├── main.go                  # Entry point - DO NOT add logic here
├── cmd/                     # CLI commands (thin layer)
│   ├── root.go              # Root command with global flags
│   ├── version.go           # Version command
│   ├── completion.go        # Shell completion
│   └── example/             # Example feature commands
│       ├── example.go       # Parent command
│       └── greet.go         # Subcommand
├── internal/                # Private application code
│   ├── handler/             # Command handlers (thin, validation)
│   │   └── greet.go         # Example handler
│   ├── service/             # Business logic (testable)
│   │   └── example/
│   │       └── service.go   # Example service
│   ├── config/              # Configuration management
│   ├── logger/              # Logging setup
│   └── model/               # Domain models and errors
├── pkg/                     # Public packages (reusable)
│   └── version/             # Version information
└── configs/                 # Configuration templates
```

## Adding a New Command

### Step 1: Create the Command File

Create `cmd/myfeature/myfeature.go`:

```go
package myfeature

import (
    "context"
    "fmt"
    "github.com/spf13/cobra"
    "github.com/yourorg/yourproject/internal/handler"
)

var Cmd = &cobra.Command{
    Use:   "myfeature",
    Short: "Brief description",
    Long:  "Longer description of what this command does",
    RunE: func(cmd *cobra.Command, args []string) error {
        h := handler.NewMyFeatureHandler()
        result, err := h.Execute(cmd.Context(), handler.MyFeatureInput{
            // Your input here
        })
        if err != nil {
            return fmt.Errorf("executing myfeature: %w", err)
        }

        fmt.Println(result.Message)
        return nil
    },
}
```

### Step 2: Register the Command

Add to `cmd/root.go`:

```go
import (
    "github.com/yourorg/yourproject/cmd/myfeature"
)

func init() {
    // ... existing code ...
    rootCmd.AddCommand(myfeature.Cmd)
}
```

### Step 3: Create the Handler

Create `internal/handler/myfeature.go`:

```go
package handler

import (
    "context"
    "fmt"
    "github.com/yourorg/yourproject/internal/service/myfeature"
)

type MyFeatureInput struct {
    Name string
}

type MyFeatureOutput struct {
    Message string
}

type myFeatureHandler struct {
    service *myfeature.Service
}

func NewMyFeatureHandler() *myFeatureHandler {
    return &myFeatureHandler{
        service: myfeature.NewService(),
    }
}

func (h *myFeatureHandler) Execute(ctx context.Context, in MyFeatureInput) (*MyFeatureOutput, error) {
    // Validation
    if in.Name == "" {
        return nil, fmt.Errorf("name is required")
    }

    // Call service
    result, err := h.service.DoSomething(ctx, in.Name)
    if err != nil {
        return nil, fmt.Errorf("doing something: %w", err)
    }

    return &MyFeatureOutput{Message: result}, nil
}
```

### Step 4: Create the Service

Create `internal/service/myfeature/service.go`:

```go
package myfeature

import (
    "context"
    "fmt"
    "log/slog"
)

type Service struct {
    // Add dependencies (repositories, clients, etc.)
}

func NewService() *Service {
    return &Service{}
}

func (s *Service) DoSomething(ctx context.Context, name string) (string, error) {
    slog.InfoContext(ctx, "doing something", "name", name)

    // Your business logic here
    return fmt.Sprintf("Processed: %s", name), nil
}
```

### Step 5: Build and Test

```bash
# Build
make build

# Test
./build/bin/yourproject myfeature

# Run tests
go test ./...
```

## Adding Flags to Commands

### Local Flags (only for this command)

```go
var myCmd = &cobra.Command{
    Use: "mycommand",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Use the flag values
        return nil
    },
}

var (
    flagName string
    flagCount int
    flagVerbose bool
)

func init() {
    // String flag
    myCmd.Flags().StringVarP(&flagName, "name", "n", "default", "description")

    // Int flag
    myCmd.Flags().IntVarP(&flagCount, "count", "c", 10, "description")

    // Bool flag
    myCmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, "description")

    // Make flag required
    myCmd.MarkFlagRequired("name")
}
```

### Persistent Flags (available to all subcommands)

Add to `cmd/root.go`:

```go
func init() {
    rootCmd.PersistentFlags().StringVar(&myFlag, "myflag", "", "description")
}
```

## Configuration

### Using Configuration Files

Create `~/.yourproject.yaml`:

```yaml
verbose: true
log_level: debug

myfeature:
  setting1: value1
  setting2: value2
```

### Reading Configuration

Update `internal/config/config.go`:

```go
type Config struct {
    Verbose    bool              `mapstructure:"verbose"`
    LogLevel   string            `mapstructure:"log_level"`
    MyFeature  MyFeatureConfig   `mapstructure:"myfeature"`
}

type MyFeatureConfig struct {
    Setting1 string `mapstructure:"setting1"`
    Setting2 string `mapstructure:"setting2"`
}
```

Update `internal/config/defaults.go`:

```go
func SetDefaults() {
    viper.SetDefault("myfeature.setting1", "default1")
    viper.SetDefault("myfeature.setting2", "default2")
}
```

Use in your code:

```go
import (
    "github.com/spf13/viper"
)

setting1 := viper.GetString("myfeature.setting1")
```

## Testing

### Unit Test Example

Create `internal/service/myfeature/service_test.go`:

```go
package myfeature

import (
    "context"
    "testing"
)

func TestService_DoSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   "test",
            want:    "Processed: test",
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
        t.Run(tt.name, func(t *testing.T) {
            s := NewService()
            got, err := s.DoSomething(context.Background(), tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("DoSomething() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if got != tt.want {
                t.Errorf("DoSomething() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Run Tests

```bash
# All tests
make test

# With coverage
make coverage

# Specific package
go test -v ./internal/service/myfeature

# With race detection
go test -race ./...
```

## Development Workflow

### Daily Development

```bash
# Format code
make fmt

# Run linters
make lint

# Fix linting issues
make lint-fix

# Run tests
make test

# Build
make build

# Run all checks
make audit
```

### Before Committing

The project uses `lefthook` for git hooks:

```bash
# Install hooks
lefthook install

# Hooks will automatically:
# - Format code (pre-commit)
# - Run linter (pre-commit)
# - Run tests (pre-push)
# - Run vet (pre-push)
```

### CI/CD

GitHub Actions workflows are configured:

- **CI** (`.github/workflows/ci.yml`): Runs on push/PR
  - Linting
  - Testing
  - Security scanning
  - Multi-platform builds

- **Release** (`.github/workflows/release.yml`): Runs on tags
  - Creates release with GoReleaser
  - Builds for multiple platforms
  - Generates changelog

## Deployment

### Docker

```bash
# Build image
docker build -f build/package/Dockerfile -t yourproject:latest .

# Run container
docker run --rm yourproject:latest version
docker run --rm yourproject:latest mycommand --help
```

### Release

```bash
# Tag a release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# GitHub Actions will automatically:
# 1. Build for all platforms
# 2. Create GitHub release
# 3. Upload binaries
```

## Best Practices

1. **Keep `cmd/` thin**: Only CLI wiring, no business logic
2. **Business logic in `service/`**: Testable, framework-independent
3. **Validate in `handler/`**: Input validation and error formatting
4. **Use structured logging**: `slog.Info("message", "key", value)`
5. **Proper error handling**: Wrap errors with context
6. **Write tests**: Table-driven tests for all logic
7. **Document flags**: Clear descriptions for all flags
8. **Use configuration**: Viper for flexible configuration

## Common Tasks

### Change Binary Name

```bash
# In Makefile
BINARY_NAME := mynewname

# In .goreleaser.yml
project_name: mynewname
binary: mynewname
```

### Add a New Dependency

```bash
go get github.com/some/package@latest
go mod tidy
```

### Generate Shell Completion

```bash
# Bash
./yourproject completion bash > /etc/bash_completion.d/yourproject

# Zsh
./yourproject completion zsh > "${fpath[1]}/_yourproject"

# Fish
./yourproject completion fish > ~/.config/fish/completions/yourproject.fish
```

## Next Steps

1. Remove the example command (optional):
   ```bash
   rm -rf cmd/example internal/handler/greet.go internal/service/example
   # Remove import and AddCommand from cmd/root.go
   ```

2. Add your own commands following the patterns above

3. Customize configuration in `internal/config/`

4. Add tests for all your code

5. Update README.md with your project details

6. Set up CI/CD with your repository

## Need Help?

- Review `GO_CLI_COMPREHENSIVE_REFERENCE.md` for detailed patterns
- Check the example command in `cmd/example/`
- Look at the existing code structure
- Follow the dependency flow: cmd → handler → service → repository

## Resources

- [Cobra Documentation](https://cobra.dev/)
- [Viper Documentation](https://github.com/spf13/viper)
- [Go Documentation](https://go.dev/doc/)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
