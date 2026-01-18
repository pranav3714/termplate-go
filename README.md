<div align="center">

# 🚀 Termplate Go

**Production-Ready Go CLI Template with Batteries Included**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/pranav3714/termplate-go)](https://goreportcard.com/report/github.com/pranav3714/termplate-go)
[![CI](https://github.com/pranav3714/termplate-go/workflows/CI/badge.svg)](https://github.com/pranav3714/termplate-go/actions)

*A modern Go CLI template with Cobra, Viper, clean architecture, and comprehensive documentation*

[Features](#-features) •
[Quick Start](#-quick-start) •
[Documentation](#-documentation) •
[Configuration](#-configuration) •
[Examples](#-examples)

</div>

---

## 🎯 Quick Navigation

> **🤖 For AI Models**: Start with [PROJECT_CONTEXT.md](PROJECT_CONTEXT.md) → [AI_GUIDE.md](AI_GUIDE.md) → [DOCUMENTATION_INDEX.md](DOCUMENTATION_INDEX.md)
>
> **👨‍💻 For Developers**: Jump to [Getting Started](#-quick-start) or [Documentation](#-documentation)

---

## ⚠️ The Honest Truth

<details>
<summary><b>Click to read the disclaimer (spoiler: it's refreshingly honest)</b></summary>

<br>

**Fair warning**: I'm not a Go expert, and this was my weekend project. You know how it is - you spend two days building something cool, and it feels like a crime to just delete it, so here we are!

**"But there's so much stuff I'll never use!"** - I hear you. This template is packed with features (API clients, database configs, multiple output formats, etc.). Think of it as a buffet - take what you need, leave what you don't. And honestly, with modern AI models like Claude, you're literally a few prompts away from removing anything you don't want. Just ask: *"Remove the database configuration"* and boom, done.

**Compatibility note**: This has been tested and works beautifully with Claude models (Sonnet/Opus/Haiku). Will it work with other AI models? Probably! I mean, it's just documentation and code structure, not rocket science. But I haven't personally tested it with GPT or others, so your mileage may vary. If you do try it with other AI models, let me know how it goes!

**TL;DR**: It works, it's documented, and it's probably overkill for your needs. But hey, better to have it and not need it than need it and have to write it on a Sunday afternoon, right? 😄

</details>

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 🛠️ Core Framework
- ⚡ **Cobra Framework** - Industry-standard CLI
- 🔧 **Viper Configuration** - Flexible config management
- 📝 **Structured Logging** - Production-ready with slog
- 🎨 **Multiple Output Formats** - JSON/YAML/Table/CSV

</td>
<td width="50%">

### 🔌 Integrations
- 🌐 **API Client Ready** - Keys, retries, rate limiting
- 📁 **File Processing** - Pattern matching, size limits
- 🗄️ **Database Support** - PostgreSQL, MySQL, SQLite
- 🎯 **Type-Safe Errors** - Domain-specific error handling

</td>
</tr>
<tr>
<td width="50%">

### 🧪 Quality & Testing
- ✅ **Comprehensive Testing** - Table-driven tests
- 🔍 **Static Analysis** - 20+ linters via golangci-lint
- 🪝 **Git Hooks** - Auto-format & lint with lefthook
- 🔒 **Security Scanning** - govulncheck integration

</td>
<td width="50%">

### 🚢 DevOps & Deployment
- 🤖 **CI/CD Ready** - GitHub Actions workflows
- 🐳 **Docker Support** - Multi-stage builds
- 📦 **Cross-Platform** - Linux, macOS, Windows
- 🎁 **Auto-Release** - GoReleaser integration

</td>
</tr>
</table>

---

## 🚀 Quick Start

### Prerequisites

- **Go** 1.22 or later
- **Make** (optional, but recommended)

### Installation

```bash
# 1. Clone the repository
git clone https://github.com/pranav3714/termplate-go.git
cd termplate-go

# 2. Install dependencies
go mod download

# 3. Build the binary
make build

# 4. Run your first command
./build/bin/termplate version
```

### First Steps

```bash
# Show all available commands
./build/bin/termplate --help

# Try the example command
./build/bin/termplate example greet --name "World"

# With different output formats
./build/bin/termplate version --output json
./build/bin/termplate version --output yaml
./build/bin/termplate version --output table

# Enable verbose logging
./build/bin/termplate -v example greet --name "User"
```

<details>
<summary><b>📋 More Usage Examples</b></summary>

```bash
# Generate shell completion
./build/bin/termplate completion bash > /etc/bash_completion.d/termplate

# Use uppercase flag
./build/bin/termplate example greet --name "User" --uppercase

# Check version with pretty JSON
./build/bin/termplate version --output json --pretty
```

</details>

---

## ⚙️ Configuration

Configuration is loaded from multiple sources (in order of priority):

| Priority | Source | Example |
|----------|--------|---------|
| 1 | Command line flag | `--config /path/to/config.yaml` |
| 2 | Environment variables | `TERMPLATE_API_KEY=xxx` |
| 3 | Home directory | `~/.termplate.yaml` |
| 4 | Current directory | `./.termplate.yaml` |

### Quick Setup

```bash
# Copy example configuration
cp configs/config.example.yaml ~/.termplate.yaml

# Edit with your settings
vim ~/.termplate.yaml
```

### Environment Variables

All config values can be set via environment variables with the `TERMPLATE_` prefix:

```bash
export TERMPLATE_API_KEY=your-api-key
export TERMPLATE_OUTPUT_FORMAT=json
export TERMPLATE_DB_USER=dbuser
export TERMPLATE_DB_PASSWORD=dbpass
```

### Configuration Examples

<details>
<summary><b>🌐 API Client Configuration</b></summary>

```yaml
api:
  base_url: https://api.github.com
  token: ${GITHUB_TOKEN}
  timeout: 60s
  retry_attempts: 3
  rate_limit: 100
  headers:
    User-Agent: termplate-go/1.0
```

</details>

<details>
<summary><b>📁 File Processing Configuration</b></summary>

```yaml
files:
  input_dir: ./data/input
  output_dir: ./data/output
  patterns: ["*.csv", "*.json"]
  max_file_size: 52428800  # 50MB
  backup_enabled: true
  create_missing_dirs: true
```

</details>

<details>
<summary><b>🗄️ Database Configuration</b></summary>

```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  database: myapp
  username: ${DB_USER}
  password: ${DB_PASSWORD}
  max_connections: 25
  connection_timeout: 30s
```

</details>

<details>
<summary><b>🎨 Output Formatting Configuration</b></summary>

```yaml
output:
  format: json        # text, json, yaml, table, csv
  pretty: true
  color: true
  table_style: unicode  # ascii, unicode, markdown
```

</details>

📖 **Full Configuration Guide**: [CONFIGURATION_GUIDE.md](docs/CONFIGURATION_GUIDE.md)

---

## 📚 Documentation

### 🤖 For AI Models (Start Here)

Perfect for AI-assisted development with Claude, GPT, and other code assistants.

| Document | Purpose |
|----------|---------|
| [**AI Guide**](AI_GUIDE.md) | Complete AI workflow guide for this codebase |
| [**Project Context**](PROJECT_CONTEXT.md) | Architecture, structure, current state |
| [**Conventions**](CONVENTIONS.md) | Coding standards, patterns, rules |
| [**Quick Reference**](QUICK_REFERENCE.md) | Fast lookups, snippets, file locations |
| [**Documentation Index**](DOCUMENTATION_INDEX.md) | Master documentation map |

### 👨‍💻 For Developers

#### Getting Started
- 🎯 [**Next Steps**](docs/NEXT_STEPS.md) - What to do now (start here!)
- 🚀 [**Getting Started**](docs/GETTING_STARTED.md) - How to add commands and customize
- ✅ [**Customization Complete**](docs/CUSTOMIZATION_COMPLETE.md) - What was customized

#### Configuration
- ⚙️ [**Configuration Guide**](docs/CONFIGURATION_GUIDE.md) - Complete configuration reference
- 📝 [**Config Example**](configs/config.example.yaml) - All available settings

#### Reference
- 📖 [**CLI Comprehensive Reference**](docs/GO_CLI_COMPREHENSIVE_REFERENCE.md) - Authoritative Go CLI patterns
- 📊 [**Project Summary**](docs/PROJECT_SUMMARY.md) - Project overview

**📑 All Documentation**: [docs/README.md](docs/README.md)

---

## 💻 Development

### Setup Development Environment

```bash
# Install development tools
make setup

# Install git hooks (pre-commit, pre-push)
lefthook install
```

### Available Make Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the binary |
| `make test` | Run all tests |
| `make coverage` | Generate coverage report |
| `make fmt` | Format code |
| `make lint` | Run linters |
| `make lint-fix` | Fix linting issues |
| `make vuln` | Check for vulnerabilities |
| `make audit` | Run all quality checks |
| `make clean` | Clean build artifacts |

### Project Structure

```
termplate/
├── 📁 cmd/                    # CLI commands (Cobra)
│   ├── root.go                # Root command with global flags
│   ├── version.go             # Version command
│   ├── completion.go          # Shell completion
│   └── example/               # Example command group
│
├── 📁 internal/               # Private application code
│   ├── config/                # Configuration (Viper)
│   ├── logger/                # Logging (slog)
│   ├── model/                 # Domain models and errors
│   ├── handler/               # Command handlers
│   ├── service/               # Business logic
│   ├── output/                # Output formatting
│   └── repository/            # Data access layer
│
├── 📁 pkg/                    # Public packages
│   └── version/               # Version information
│
├── 📁 configs/                # Configuration templates
├── 📁 docs/                   # Documentation
├── 📁 test/                   # Integration & E2E tests
└── 📁 build/                  # Build configurations
```

---

## 🎓 Examples

### Adding Your First Command

<details>
<summary><b>Click to see a complete example</b></summary>

<br>

**Step 1**: Create command file `cmd/mycommand/mycommand.go`

```go
package mycommand

import (
    "github.com/spf13/cobra"
    "github.com/pranav3714/termplate-go/internal/handler"
)

var Cmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Description of your command",
    Long: `A longer description that spans multiple lines and
    explains what your command does in detail.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        h := handler.NewMyHandler()
        return h.Execute(cmd.Context(), handler.MyInput{
            Name: name,
        })
    },
}

var name string

func init() {
    Cmd.Flags().StringVarP(&name, "name", "n", "", "Your name")
    Cmd.MarkFlagRequired("name")
}
```

**Step 2**: Register command in `cmd/root.go`

```go
import "github.com/pranav3714/termplate-go/cmd/mycommand"

func init() {
    rootCmd.AddCommand(mycommand.Cmd)
}
```

**Step 3**: Run your command

```bash
./build/bin/termplate mycommand --name "World"
```

📖 **Detailed Guide**: [Getting Started](docs/GETTING_STARTED.md)

</details>

---

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Generate HTML coverage report
make coverage
open coverage.html
```

---

## 🏗️ Building

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
docker build -f build/package/Dockerfile -t termplate:latest .

# Run in container
docker run --rm termplate:latest version
```

### Cross-Platform Releases

This project uses [GoReleaser](https://goreleaser.com/) for automated releases.

```bash
# Create and push a tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# GitHub Actions automatically builds and releases for:
# - Linux (amd64, arm64)
# - macOS (amd64, arm64)
# - Windows (amd64, arm64)
```

---

## 📦 Releasing

Quick release automation using GitHub Actions and GoReleaser.

### Quick Commands

```bash
# Automated (recommended)
make release-prepare

# Preview first
make release-dry-run

# Auto-increment patch
make release-patch
```

### Release Flow

1. **Update code** and commit changes
2. **Run** `make release-prepare`
3. **Script automates**:
   - CHANGELOG.md updates
   - Version links
   - Git tag creation
4. **GitHub Actions automatically**:
   - Builds 6 platforms (Linux, macOS, Windows - amd64, arm64)
   - Creates release
   - Uploads artifacts

### Complete Guide

📖 **[RELEASE_RULEBOOK.md](RELEASE_RULEBOOK.md)** - Complete release documentation including:
- Automated and manual processes
- Troubleshooting guide
- Rollback procedures
- Version numbering guidelines
- AI-assisted releases

---

## 🐚 Shell Completion

Enable shell completion for a better CLI experience:

```bash
# Bash
termplate completion bash > /etc/bash_completion.d/termplate

# Zsh
termplate completion zsh > "${fpath[1]}/_termplate"

# Fish
termplate completion fish > ~/.config/fish/completions/termplate.fish

# PowerShell
termplate completion powershell | Out-String | Invoke-Expression
```

---

## 🏛️ Architecture

This CLI follows **clean architecture** principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────┐
│                        main.go                          │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│                    cmd/ (Cobra)                         │
│     • CLI wiring only, no business logic               │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│              internal/handler/                          │
│     • Input validation and orchestration               │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│              internal/service/                          │
│     • Business logic (testable, framework-free)        │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│            internal/repository/                         │
│     • Data access abstraction                          │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│               internal/model/                           │
│     • Domain models, errors, interfaces                │
└─────────────────────────────────────────────────────────┘
```

**Dependencies flow downward only** → No circular dependencies

---

## ✅ Best Practices

This template follows industry-standard Go best practices:

<table>
<tr>
<td width="50%">

**Project Structure**
- ✅ Standard Project Layout
- ✅ Clear separation of concerns
- ✅ Dependency injection
- ✅ Interface-based design

**Code Quality**
- ✅ Comprehensive error handling
- ✅ Structured logging
- ✅ Table-driven tests
- ✅ Static analysis (20+ linters)

</td>
<td width="50%">

**DevOps**
- ✅ Git hooks (lefthook)
- ✅ CI/CD (GitHub Actions)
- ✅ Automated releases (GoReleaser)
- ✅ Docker support

**Developer Experience**
- ✅ Shell completion
- ✅ Comprehensive documentation
- ✅ AI-optimized (Claude/GPT)
- ✅ Example commands

</td>
</tr>
</table>

---

## 📖 Resources

Recommended reading and resources:

- [Cobra Documentation](https://cobra.dev/) - CLI framework
- [Viper Documentation](https://github.com/spf13/viper) - Configuration management
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout) - Project structure
- [Effective Go](https://go.dev/doc/effective_go) - Go best practices
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) - Code style

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 💬 Support

- 📚 Check the [documentation](docs/) for detailed guides
- 🐛 Report bugs via [GitHub Issues](https://github.com/pranav3714/termplate-go/issues)
- 💡 Request features via [GitHub Issues](https://github.com/pranav3714/termplate-go/issues)
- 🔒 Report security issues via [Security Policy](SECURITY.md)

---

<div align="center">

**Made with ❤️ on a weekend**

*Because good templates shouldn't be reinvented every time*

[⬆ Back to Top](#-termplate-go)

</div>
