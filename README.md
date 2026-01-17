# ever-so-powerful-go

A production-ready Go CLI tool built with best practices and a strong foundation. Designed for developers who need API interaction, file processing, database operations, and flexible output formatting.

## Features

- **Cobra Framework**: Industry-standard CLI with nested commands, auto-help, and shell completion
- **Viper Configuration**: Flexible configuration with config files, environment variables, and defaults
- **Structured Logging**: Production-ready logging with slog (JSON/text output)
- **Multiple Output Formats**: JSON, YAML, Table (ASCII/Unicode/Markdown), CSV
- **API Client Configuration**: Keys, tokens, retries, rate limiting, custom headers
- **File Processing**: Pattern matching, size limits, backup options, permission control
- **Database Support**: PostgreSQL, MySQL, SQLite with connection pooling
- **Type-Safe Error Handling**: Domain-specific errors with proper wrapping
- **Comprehensive Testing**: Table-driven tests, mocks, and test helpers
- **Static Analysis**: golangci-lint with 20+ linters
- **CI/CD Ready**: GitHub Actions workflows for testing and releasing
- **Cross-Platform**: Build for Linux, macOS, and Windows
- **Docker Support**: Multi-stage Dockerfile for minimal container images
- **Git Hooks**: Automatic formatting and linting with lefthook

## Quick Start

### Prerequisites

- Go 1.22 or later
- Make (optional, for using Makefile commands)

### Installation

```bash
# Clone the repository
git clone https://github.com/blacksilver/ever-so-powerful.git
cd ever-so-powerful

# Install dependencies
go mod download

# Build the binary
make build

# Or build manually
go build -o build/bin/ever-so-powerful-go .
```

### Usage

```bash
# Show help
./build/bin/ever-so-powerful-go --help

# Show version
./build/bin/ever-so-powerful-go version

# Show version in JSON format
./build/bin/ever-so-powerful-go version --output json

# Run example command
./build/bin/ever-so-powerful-go example greet --name "World"

# With verbose logging
./build/bin/ever-so-powerful-go -v example greet --name "User" --uppercase

# Generate shell completion
./build/bin/ever-so-powerful-go completion bash > /etc/bash_completion.d/ever-so-powerful-go
```

## Configuration

Configuration can be loaded from:

1. **Command line flag**: `--config /path/to/config.yaml`
2. **Home directory**: `~/.ever-so-powerful-go.yaml`
3. **Current directory**: `./.ever-so-powerful-go.yaml`
4. **Environment variables**: `EVER_SO_POWERFUL_GO_*`

### Quick Setup

```bash
# Copy example configuration
cp configs/config.example.yaml ~/.ever-so-powerful-go.yaml

# Edit with your settings
vim ~/.ever-so-powerful-go.yaml

# Or use environment variables
export EVER_SO_POWERFUL_GO_API_KEY=your-api-key
export EVER_SO_POWERFUL_GO_OUTPUT_FORMAT=json
export EVER_SO_POWERFUL_GO_DB_USER=dbuser
export EVER_SO_POWERFUL_GO_DB_PASSWORD=dbpass
```

See [Configuration Guide](docs/CONFIGURATION_GUIDE.md) for detailed documentation.

## Documentation

- **[Next Steps](docs/NEXT_STEPS.md)** - What to do now (start here!)
- **[Getting Started](docs/GETTING_STARTED.md)** - How to add commands and customize
- **[Configuration Guide](docs/CONFIGURATION_GUIDE.md)** - Complete configuration reference
- **[Customization Complete](docs/CUSTOMIZATION_COMPLETE.md)** - What was customized
- **[Project Summary](docs/PROJECT_SUMMARY.md)** - Project overview
- **[CLI Reference](docs/GO_CLI_COMPREHENSIVE_REFERENCE.md)** - Complete Go CLI reference

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

## Project Structure

```
ever-so-powerful/
├── cmd/                          # CLI commands (Cobra)
│   ├── root.go                   # Root command with global flags
│   ├── version.go                # Version command
│   ├── completion.go             # Shell completion
│   └── example/                  # Example command group
├── internal/                     # Private application code
│   ├── config/                   # Configuration (Viper)
│   ├── logger/                   # Logging (slog)
│   ├── model/                    # Domain models and errors
│   ├── handler/                  # Command handlers
│   ├── service/                  # Business logic
│   ├── output/                   # Output formatting
│   └── repository/               # Data access layer
├── pkg/                          # Public packages
│   └── version/                  # Version information
├── configs/                      # Configuration templates
├── docs/                         # Documentation
├── test/                         # Integration & E2E tests
└── build/                        # Build configurations
```

## Adding Your First Command

See the detailed guide in [Getting Started](docs/GETTING_STARTED.md).

Quick example:

```go
// cmd/mycommand/mycommand.go
package mycommand

import (
    "github.com/spf13/cobra"
    "github.com/blacksilver/ever-so-powerful/internal/handler"
)

var Cmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Description of your command",
    RunE: func(cmd *cobra.Command, _ []string) error {
        h := handler.NewMyHandler()
        return h.Execute(cmd.Context(), handler.MyInput{})
    },
}
```

Register in `cmd/root.go`:

```go
import "github.com/blacksilver/ever-so-powerful/cmd/mycommand"

func init() {
    rootCmd.AddCommand(mycommand.Cmd)
}
```

## Configuration Examples

### API Client

```yaml
api:
  base_url: https://api.github.com
  token: ${GITHUB_TOKEN}
  timeout: 60s
  retry_attempts: 3
```

### File Processing

```yaml
files:
  input_dir: ./data/input
  output_dir: ./data/output
  patterns: ["*.csv", "*.json"]
  max_file_size: 52428800  # 50MB
```

### Database

```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  database: myapp
  username: ${DB_USER}
  password: ${DB_PASSWORD}
```

### Output Formatting

```yaml
output:
  format: json        # text, json, yaml, table, csv
  pretty: true
  color: true
  table_style: unicode  # ascii, unicode, markdown
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Generate coverage report
make coverage
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
docker build -f build/package/Dockerfile -t ever-so-powerful-go:latest .

# Run Docker container
docker run --rm ever-so-powerful-go:latest version
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
ever-so-powerful-go completion bash > /etc/bash_completion.d/ever-so-powerful-go

# Zsh
ever-so-powerful-go completion zsh > "${fpath[1]}/_ever-so-powerful-go"

# Fish
ever-so-powerful-go completion fish > ~/.config/fish/completions/ever-so-powerful-go.fish

# PowerShell
ever-so-powerful-go completion powershell | Out-String | Invoke-Expression
```

## Architecture

This CLI follows clean architecture principles:

- **cmd/**: Thin CLI layer - only wiring, no business logic
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

## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Resources

- [Cobra Documentation](https://cobra.dev/)
- [Viper Documentation](https://github.com/spf13/viper)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Effective Go](https://go.dev/doc/effective_go)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

## Support

For detailed guides and examples, check the [documentation](docs/).
