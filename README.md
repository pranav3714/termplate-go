# mycli

A production-ready Go CLI tool built with best practices and a strong foundation.

## Features

- **Cobra Framework**: Industry-standard CLI framework with nested commands, auto-help, and shell completion
- **Viper Configuration**: Flexible configuration with support for config files, environment variables, and defaults
- **Structured Logging**: Production-ready logging with slog (JSON in production, text in development)
- **Type-Safe Error Handling**: Domain-specific errors with proper wrapping and unwrapping
- **Comprehensive Testing**: Table-driven tests, mocks, and test helpers
- **Static Analysis**: golangci-lint with extensive linter suite
- **CI/CD Ready**: GitHub Actions workflows for testing and releasing
- **Cross-Platform**: Build for Linux, macOS, and Windows
- **Docker Support**: Multi-stage Dockerfile for minimal container images
- **Release Automation**: GoReleaser configuration for automated releases

## Quick Start

### Prerequisites

- Go 1.22 or later
- Make (optional, for using Makefile commands)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourorg/mycli.git
cd mycli

# Install dependencies
go mod download

# Build the binary
make build

# Or build manually
go build -o build/bin/mycli .
```

### Usage

```bash
# Show help
./build/bin/mycli --help

# Show version
./build/bin/mycli version

# Show version in JSON format
./build/bin/mycli version --output json

# Generate shell completion
./build/bin/mycli completion bash > /etc/bash_completion.d/mycli
```

## Development

### Setup Development Environment

```bash
# Install development tools (golangci-lint, govulncheck, goimports, lefthook)
make setup

# Install git hooks
lefthook install
```

### Common Tasks

```bash
# Build the binary
make build

# Run tests
make test

# Run tests with coverage
make coverage

# Format code
make fmt

# Run linters
make lint

# Fix linting issues
make lint-fix

# Check for vulnerabilities
make vuln

# Run all quality checks
make audit

# Clean build artifacts
make clean
```

### Project Structure

```
mycli/
├── main.go                       # Entry point (minimal code)
├── cmd/                          # CLI commands (Cobra)
│   ├── root.go                   # Root command, global flags
│   ├── version.go                # Version command
│   └── completion.go             # Shell completion command
├── internal/                     # Private application code
│   ├── config/                   # Configuration (Viper)
│   ├── logger/                   # Logging (slog)
│   ├── model/                    # Domain models and errors
│   ├── handler/                  # Command handlers
│   ├── service/                  # Business logic
│   └── repository/               # Data access layer
├── pkg/                          # Public packages
│   └── version/                  # Version information
├── configs/                      # Configuration templates
│   └── config.example.yaml       # Example configuration
├── test/                         # Integration & E2E tests
├── build/                        # Build configurations
│   └── package/                  # Docker files
└── scripts/                      # Development scripts
```

## Configuration

The CLI tool supports multiple configuration sources with the following priority (highest to lowest):

1. Explicit flags
2. Environment variables (prefixed with `MYCLI_`)
3. Configuration file (`.mycli.yaml` in home directory or current directory)
4. Default values

### Configuration File

Create a configuration file at `~/.mycli.yaml` or `./.mycli.yaml`:

```yaml
verbose: false
log_level: info
output: text

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

### Environment Variables

```bash
export MYCLI_VERBOSE=true
export MYCLI_LOG_LEVEL=debug
export MYCLI_API_KEY=your-api-key
```

## Adding New Commands

To add a new command, follow this pattern:

1. Create a new file in `cmd/` (e.g., `cmd/mycommand.go`)
2. Define the command structure
3. Add the command to `rootCmd` in `cmd/root.go`

Example:

```go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Brief description",
    Long:  "Longer description",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Hello from mycommand!")
        return nil
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific test
go test -v ./internal/config -run TestLoad
```

## Building

### Local Build

```bash
# Build for current platform
make build

# Build for all platforms
make build-all
```

### Docker Build

```bash
# Build Docker image
docker build -f build/package/Dockerfile -t mycli:latest .

# Run Docker container
docker run --rm mycli:latest version
```

## Releasing

This project uses GoReleaser for automated releases.

```bash
# Create a git tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# GoReleaser will automatically build and release via GitHub Actions
# Or run locally:
goreleaser release --clean
```

## Shell Completion

Generate shell completion scripts:

```bash
# Bash
mycli completion bash > /etc/bash_completion.d/mycli
source ~/.bashrc

# Zsh
mycli completion zsh > "${fpath[1]}/_mycli"
source ~/.zshrc

# Fish
mycli completion fish > ~/.config/fish/completions/mycli.fish
source ~/.config/fish/config.fish

# PowerShell
mycli completion powershell | Out-String | Invoke-Expression
```

## License

[Your License Here]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Customization Guide

To customize this template for your project:

1. **Update Module Path**: Change `github.com/yourorg/mycli` to your actual module path
   - In `go.mod`
   - In all import statements
   - In `Makefile`
   - In `.goreleaser.yml`

2. **Update Binary Name**: Change `mycli` to your desired binary name
   - In `Makefile` (`BINARY_NAME`)
   - In `README.md`
   - In `.goreleaser.yml`

3. **Update Description**: Edit the descriptions in
   - `cmd/root.go` (`Short` and `Long` fields)
   - `README.md`

4. **Configure Your Domain**: Update the configuration structures in
   - `internal/config/config.go`
   - `internal/config/defaults.go`
   - `configs/config.example.yaml`

5. **Add Your Commands**: Create new command files in `cmd/` directory

6. **Implement Business Logic**: Add your services in `internal/service/`

7. **Add Tests**: Write tests for all your code

## Architecture

This CLI tool follows clean architecture principles:

- **cmd/**: Thin CLI layer, only wiring
- **internal/handler/**: Input validation and command handling
- **internal/service/**: Business logic (testable, framework-independent)
- **internal/repository/**: Data access abstraction
- **internal/model/**: Domain models and errors
- **pkg/**: Reusable, public packages

Dependencies flow downward:
```
main.go → cmd/ → handler/ → service/ → repository/ → model/
```

## Best Practices

This template follows Go best practices:

- ✅ Standard Project Layout
- ✅ Clear separation of concerns
- ✅ Dependency injection
- ✅ Interface-based design
- ✅ Comprehensive error handling
- ✅ Structured logging
- ✅ Table-driven tests
- ✅ Static analysis with golangci-lint
- ✅ Git hooks with lefthook
- ✅ CI/CD with GitHub Actions
- ✅ Automated releases with GoReleaser
- ✅ Docker support
- ✅ Shell completion

## References

- [Cobra Documentation](https://cobra.dev/)
- [Viper Documentation](https://github.com/spf13/viper)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Effective Go](https://go.dev/doc/effective_go)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
