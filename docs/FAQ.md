# Frequently Asked Questions (FAQ)

> **Document Type**: Reference Guide
> **Purpose**: Quick answers to common questions about Termplate Go
> **Keywords**: faq, questions, help, troubleshoot, common issues, how to
> **Related**: GETTING_STARTED.md, DEBUG_GUIDE.md, QUICK_REFERENCE.md

Quick answers to common questions about developing with Termplate Go.

## Table of Contents

- [Getting Started](#getting-started)
- [Development](#development)
- [Testing](#testing)
- [Configuration](#configuration)
- [Building & Deployment](#building--deployment)
- [Debugging](#debugging)
- [Performance](#performance)
- [Contributing](#contributing)

---

## Getting Started

### Q: How do I start using this template?

**A:** Follow these steps:

```bash
# 1. Clone or use as template
git clone <your-repo>

# 2. Update module name
sed -i 's|github.com/yourorg/mycli|github.com/yourusername/yourproject|g' go.mod

# 3. Install dependencies
go mod tidy

# 4. Build
make build

# 5. Test
./build/bin/mycli --help
```

See **GETTING_STARTED.md** for detailed instructions.

### Q: What Go version is required?

**A:** Go 1.23 or later. Check with:

```bash
go version
```

Update go.mod if needed:

```bash
go mod edit -go=1.24
```

### Q: What's included in this template?

**A:** Core features:
- CLI framework (Cobra + Viper)
- Structured logging (slog)
- Configuration management
- Output formatting (JSON, YAML, Table, CSV)
- Testing infrastructure
- CI/CD (GitHub Actions)
- Docker support
- Release automation

See **PROJECT_CONTEXT.md** for complete overview.

### Q: Can I use this for commercial projects?

**A:** Yes! This is a starter template for your projects. The license terms depend on what license you choose for your project.

---

## Development

### Q: How do I add a new command?

**A:** Follow the clean architecture pattern:

1. **Create command file** (`cmd/myfeature/myfeature.go`)
2. **Register command** (add to `cmd/root.go`)
3. **Create handler** (`internal/handler/myfeature.go`)
4. **Create service** (`internal/service/myfeature/service.go`)
5. **Add tests**

Example:

```bash
# Quick way
cp -r cmd/example cmd/myfeature
# Then customize files
```

See **GETTING_STARTED.md** → "Adding a New Command" for details.

### Q: Where does business logic go?

**A:** Follow this separation:

- `cmd/` - CLI wiring only, no logic
- `internal/handler/` - Input validation, error formatting
- `internal/service/` - **Business logic goes here**
- `internal/model/` - Domain models, errors

**Rule:** If it can be tested without Cobra, it belongs in `service/`.

### Q: How do I add configuration options?

**A:** Three steps:

1. **Add to struct** (`internal/config/config.go`):
   ```go
   type Config struct {
       MyFeature MyFeatureConfig `mapstructure:"myfeature"`
   }
   ```

2. **Set defaults** (`internal/config/defaults.go`):
   ```go
   viper.SetDefault("myfeature.setting", "default-value")
   ```

3. **Use in code**:
   ```go
   value := viper.GetString("myfeature.setting")
   ```

See **CONFIGURATION_GUIDE.md** for all options.

### Q: What's the code style standard?

**A:** Follow **CONVENTIONS.md**:

- `gofmt` for formatting
- Package-level errors (`var ErrNotFound = errors.New(...)`)
- Context as first parameter
- Structured logging with slog
- Table-driven tests

Enforced by:
```bash
make lint        # Check
make lint-fix    # Auto-fix
```

### Q: How do I update dependencies?

**A:** Use these commands:

```bash
# Update all (within constraints)
go get -u ./...
go mod tidy

# Update specific package
go get github.com/some/package@latest

# Check for outdated
go list -u -m all

# Verify
make test
make lint
```

---

## Testing

### Q: How do I run tests?

**A:** Multiple options:

```bash
make test           # All tests
make coverage       # With coverage report

# Specific package
go test -v ./internal/service/myservice

# Specific test
go test -v ./internal/service/myservice -run TestMyFunction

# With race detector
go test -race ./...

# Integration tests
go test -v -tags=integration ./...
```

### Q: What should I test?

**A:** Focus on:
- ✅ Business logic in services
- ✅ Input validation in handlers
- ✅ Error conditions
- ✅ Edge cases

Skip:
- ❌ Cobra command setup (framework code)
- ❌ Simple getters/setters
- ❌ Third-party libraries

### Q: How do I mock dependencies?

**A:** Use interface-based mocking:

```go
// Define interface
type Repository interface {
    Get(ctx context.Context, id string) (*Model, error)
}

// Create mock
type MockRepository struct {
    GetFunc func(ctx context.Context, id string) (*Model, error)
}

func (m *MockRepository) Get(ctx context.Context, id string) (*Model, error) {
    if m.GetFunc != nil {
        return m.GetFunc(ctx, id)
    }
    return nil, errors.New("not implemented")
}

// Use in test
mock := &MockRepository{
    GetFunc: func(ctx context.Context, id string) (*Model, error) {
        return &Model{ID: id}, nil
    },
}
```

See **TESTING_PATTERNS.md** for comprehensive examples.

### Q: What's the target coverage?

**A:** Recommended:
- Overall: 70-80%
- Business logic (services): 90%+
- Handlers: 80%+
- Commands: 50-60%

Check with:
```bash
make coverage
go tool cover -func=coverage.out | grep total
```

---

## Configuration

### Q: Where do config files go?

**A:** Multiple locations (in priority order):

1. `./config.yaml` (current directory)
2. `~/.mycli.yaml` (home directory)
3. `/etc/mycli/config.yaml` (system-wide)
4. `--config=/path/to/config.yaml` (explicit flag)

### Q: How do I use environment variables?

**A:** Prefix with `TERMPLATE_`:

```bash
export TERMPLATE_LOG_LEVEL=debug
export TERMPLATE_VERBOSE=true
export TERMPLATE_API_KEY=secret

./mycli mycommand
```

All config values can be overridden via env vars.

### Q: What config formats are supported?

**A:** Viper supports:
- YAML (`.yaml`, `.yml`)
- JSON (`.json`)
- TOML (`.toml`)
- ENV (environment variables)
- Command-line flags

Recommended: YAML for readability.

### Q: How do flags override config?

**A:** Priority order (highest to lowest):

1. Command-line flags
2. Environment variables
3. Config file
4. Default values

Example:
```bash
# Config file: log_level: info
# Override with flag:
./mycli --log-level=debug mycommand

# Override with env:
TERMPLATE_LOG_LEVEL=debug ./mycli mycommand
```

---

## Building & Deployment

### Q: How do I build for production?

**A:** Use the Makefile:

```bash
make build          # Current platform
make build-all      # All platforms (Linux, macOS, Windows)

# Manual
CGO_ENABLED=0 go build -ldflags="-w -s" -o mycli ./main.go
```

Artifacts in `build/bin/`.

### Q: How do I create a release?

**A:** Use automated release system:

```bash
# Interactive
make release-prepare

# Preview first
make release-dry-run

# Auto-increment patch
make release-patch
```

See **RELEASE_RULEBOOK.md** for complete guide.

### Q: How does CI/CD work?

**A:** GitHub Actions workflows:

- **CI** (`.github/workflows/ci.yml`):
  - Runs on push/PR
  - Linting, testing, security scan
  - Multi-platform builds

- **Release** (`.github/workflows/release.yml`):
  - Runs on git tags (`v*`)
  - GoReleaser builds for 6 platforms
  - Creates GitHub release
  - Uploads binaries

Trigger release:
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### Q: How do I use Docker?

**A:** Multiple options:

```bash
# Build image
docker build -f build/package/Dockerfile -t mycli:latest .

# Run
docker run --rm mycli:latest version
docker run --rm mycli:latest mycommand

# Development
docker-compose up -d
docker-compose exec app make build
```

See **DOCKER_DEVELOPMENT.md** for comprehensive guide.

### Q: What platforms are supported?

**A:** GoReleaser builds for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

Total: 6 platform/architecture combinations.

---

## Debugging

### Q: How do I enable debug logging?

**A:** Multiple ways:

```bash
# Flag
./mycli -v mycommand
./mycli --log-level=debug mycommand

# Environment
export TERMPLATE_LOG_LEVEL=debug
./mycli mycommand

# Config file
# config.yaml:
log_level: debug
```

### Q: How do I debug with Delve?

**A:** Install and use:

```bash
# Install
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug ./main.go -- mycommand --flag value

# Set breakpoint
(dlv) break main.main
(dlv) break handler.go:42
(dlv) continue

# Inspect
(dlv) print myVar
(dlv) locals
(dlv) stack
```

See **DEBUG_GUIDE.md** for comprehensive debugging strategies.

### Q: Tests pass but CLI fails?

**A:** Common causes:

1. **Configuration issue:**
   ```bash
   ./mycli --config=/path/to/config.yaml -v mycommand
   ```

2. **Context not propagated:**
   - Ensure `cmd.Context()` passed through handlers

3. **Missing initialization:**
   - Check `init()` functions in cmd/ packages

4. **Different working directory:**
   - Use absolute paths or check `os.Getwd()`

### Q: How do I profile performance?

**A:** Add profiling flags:

```go
// cmd/mycommand/mycommand.go
var cpuProfile string

func init() {
    Cmd.Flags().StringVar(&cpuProfile, "cpuprofile", "", "CPU profile")
}
```

Then:
```bash
./mycli mycommand --cpuprofile=cpu.prof
go tool pprof cpu.prof
```

See **PERFORMANCE_GUIDE.md** for detailed profiling.

---

## Performance

### Q: How do I benchmark my code?

**A:** Write benchmark tests:

```go
// service_bench_test.go
func BenchmarkProcess(b *testing.B) {
    s := NewService()
    input := "test"

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        s.Process(context.Background(), input)
    }
}
```

Run:
```bash
go test -bench=. -benchmem ./internal/service/myservice
```

### Q: How do I reduce memory allocations?

**A:** Common optimizations:

1. **Pre-allocate slices:**
   ```go
   result := make([]string, 0, len(items))
   ```

2. **Use sync.Pool for buffers:**
   ```go
   var bufferPool = sync.Pool{
       New: func() interface{} { return new(bytes.Buffer) },
   }
   ```

3. **Use strings.Builder:**
   ```go
   var sb strings.Builder
   sb.Grow(estimatedSize)
   sb.WriteString(str)
   ```

See **PERFORMANCE_GUIDE.md** for comprehensive optimization techniques.

### Q: How do I check for memory leaks?

**A:** Use profiling:

```bash
# Run with memory profile
./mycli mycommand --memprofile=mem.prof

# Analyze
go tool pprof -alloc_space mem.prof
(pprof) top
(pprof) list FunctionName

# Check goroutines
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

### Q: How do I optimize Docker images?

**A:** Follow best practices:

1. **Use multi-stage builds** (already in template)
2. **Use .dockerignore**
3. **Minimize layers**
4. **Use specific tags** (not `latest`)
5. **Build with BuildKit:**
   ```bash
   DOCKER_BUILDKIT=1 docker build -f build/package/Dockerfile .
   ```

Result: ~5-10MB images.

---

## Contributing

### Q: What checks run before commit?

**A:** Lefthook git hooks:

- **Pre-commit:**
  - `go fmt`
  - `golangci-lint run`

- **Pre-push:**
  - `go test ./...`
  - `go vet ./...`

Install hooks:
```bash
lefthook install
```

### Q: How do I run all quality checks?

**A:** Use the audit target:

```bash
make audit
```

Runs:
- `go fmt`
- `go vet`
- `golangci-lint`
- `go test`
- `go mod verify`

### Q: What linters are configured?

**A:** 20+ linters in `.golangci.yml`:

Core linters:
- `errcheck` - Unchecked errors
- `gosimple` - Simplifications
- `govet` - Go vet issues
- `ineffassign` - Ineffective assignments
- `staticcheck` - Static analysis
- `unused` - Unused code

Run:
```bash
make lint        # Check
make lint-fix    # Auto-fix where possible
```

### Q: How do I contribute changes?

**A:** Standard workflow:

1. **Create branch:**
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make changes and test:**
   ```bash
   make audit
   ```

3. **Commit:**
   ```bash
   git commit -m "feat: add new feature"
   ```

4. **Push:**
   ```bash
   git push origin feature/my-feature
   ```

5. **Create pull request**

---

## Common Error Messages

### Error: `module not found`

**Cause:** Dependencies not installed.

**Fix:**
```bash
go mod download
go mod tidy
```

### Error: `command not found: mycli`

**Cause:** Binary not built or not in PATH.

**Fix:**
```bash
make build
export PATH=$PATH:$(pwd)/build/bin
# Or use full path:
./build/bin/mycli mycommand
```

### Error: `config file not found`

**Cause:** Config file in wrong location.

**Fix:**
```bash
# Check locations
ls -la ./.mycli.yaml
ls -la ~/.mycli.yaml

# Or specify explicitly
./mycli --config=/path/to/config.yaml mycommand
```

### Error: `permission denied`

**Cause:** Binary not executable.

**Fix:**
```bash
chmod +x build/bin/mycli
```

### Error: `import cycle not allowed`

**Cause:** Circular dependency between packages.

**Fix:**
- Extract shared code to separate package
- Move toward handler → service → model flow
- Avoid service ↔ service dependencies

### Error: `undefined: SomeFunction`

**Cause:** Function/type not exported (lowercase first letter).

**Fix:**
```go
// Bad
func myFunction() {}

// Good
func MyFunction() {}
```

---

## Quick Reference

### Essential Commands

```bash
# Development
make build              # Build binary
make test               # Run tests
make lint               # Run linters
make audit              # All checks

# Docker
docker build -f build/package/Dockerfile -t mycli:latest .
docker run --rm mycli:latest version

# Release
make release-prepare    # Create release
make release-dry-run    # Preview

# Help
make help               # Show all targets
./mycli --help          # CLI help
```

### File Locations

```
Essential files:
├── cmd/                        # CLI commands
├── internal/
│   ├── handler/                # Command handlers
│   ├── service/                # Business logic
│   ├── config/                 # Configuration
│   └── model/                  # Domain models
├── docs/                       # Documentation
├── Makefile                    # Build automation
├── .golangci.yml               # Linter config
└── lefthook.yml                # Git hooks
```

### Documentation Map

- **Getting started** → GETTING_STARTED.md
- **AI workflows** → AI_GUIDE.md
- **Code standards** → CONVENTIONS.md
- **Quick lookups** → QUICK_REFERENCE.md
- **Debugging** → DEBUG_GUIDE.md
- **Testing** → TESTING_PATTERNS.md
- **Docker** → DOCKER_DEVELOPMENT.md
- **Performance** → PERFORMANCE_GUIDE.md
- **Release** → RELEASE_RULEBOOK.md
- **All docs** → DOCUMENTATION_INDEX.md

---

## Still Need Help?

1. **Check documentation:**
   - See **DOCUMENTATION_INDEX.md** for complete guide
   - Use search in docs/ directory

2. **Review examples:**
   - Look at `cmd/example/` for working command
   - Check tests for usage patterns

3. **Enable debug logging:**
   ```bash
   ./mycli -v mycommand --log-level=debug
   ```

4. **Ask AI:**
   - "How do I add a new command in Termplate Go?"
   - "Show me testing patterns in Termplate Go"
   - AI can read **AI_GUIDE.md** for context

5. **Check GitHub issues:**
   - Search existing issues
   - Create new issue with:
     - Go version
     - OS/platform
     - Command that failed
     - Full error message
     - Steps to reproduce
