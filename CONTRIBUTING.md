# Contributing to Termplate Go

First off, thanks for taking the time to contribute! 🎉

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Submitting Changes](#submitting-changes)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)

## Code of Conduct

This project adheres to a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the issue
- **Expected behavior** vs **actual behavior**
- **Environment details** (OS, Go version, etc.)
- **Code samples** or **error messages** if applicable

**Use the bug report template** when creating an issue.

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide detailed description** of the suggested enhancement
- **Explain why this enhancement would be useful**
- **List similar features** in other tools if applicable

**Use the feature request template** when creating an issue.

### Pull Requests

- Fill in the pull request template
- Follow the [coding standards](#coding-standards)
- Include tests for new features
- Update documentation as needed
- Ensure all CI checks pass

## Development Setup

### Prerequisites

- Go 1.22 or later
- Make (optional but recommended)
- Git

### Setup Steps

```bash
# Clone the repository
git clone https://github.com/pranav3714/termplate-go.git
cd termplate-go

# Install development tools
make setup

# Install git hooks (automatic formatting and linting)
lefthook install

# Build the binary
make build

# Run tests
make test
```

## Coding Standards

This project follows strict coding standards to maintain quality and consistency.

### Required Reading

Before contributing code, please read:

1. **[CONVENTIONS.md](CONVENTIONS.md)** - Coding standards and patterns
2. **[PROJECT_CONTEXT.md](PROJECT_CONTEXT.md)** - Project architecture
3. **[cmd/example/greet.go](cmd/example/greet.go)** - Working example

### Architecture Rules

This project uses clean architecture with strict layer separation:

```
cmd/              → CLI wiring ONLY (no business logic)
internal/handler/ → Validation + I/O formatting
internal/service/ → Business logic (testable, framework-independent)
internal/repository/ → Data access
internal/model/   → Domain entities (no dependencies)
```

**Key Rules**:
- ✅ Dependencies only flow downward
- ✅ No circular imports
- ✅ Business logic stays in `service/` layer
- ❌ Never put business logic in `cmd/`

### Code Quality Requirements

All contributions must:

1. **Pass linting**: Run `make lint` (20+ linters)
2. **Pass tests**: Run `make test`
3. **Pass vet**: Run `make vet`
4. **Be formatted**: Run `make fmt` (done automatically by git hooks)
5. **Have no vulnerabilities**: Run `make vuln`

**Quick check**: Run `make audit` to run all quality checks at once.

### Error Handling

Always wrap errors with context:

```go
// ✅ GOOD
if err != nil {
    return fmt.Errorf("processing file: %w", err)
}

// ❌ BAD
if err != nil {
    return err
}
```

### Logging

Use structured logging with `slog`:

```go
// ✅ GOOD
slog.InfoContext(ctx, "processing item", "id", itemID, "count", count)

// ❌ BAD
fmt.Printf("Processing item %s with count %d\n", itemID, count)
```

### Naming Conventions

- **Packages**: lowercase, single word (e.g., `mypackage`)
- **Files**: lowercase with underscores (e.g., `my_handler.go`)
- **Exported types**: PascalCase (e.g., `MyHandler`)
- **Unexported types**: camelCase (e.g., `myInternalType`)
- **Functions**: PascalCase (exported) or camelCase (unexported)

## Submitting Changes

### Before Submitting

1. **Read the documentation**: Understand the architecture and patterns
2. **Check existing issues**: See if your change is already discussed
3. **Create an issue**: Discuss large changes before implementing
4. **Follow conventions**: Match existing code style

### Commit Message Format

```
<type>: <short description>

<optional longer description>

<optional footer>
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `chore`: Build process, dependencies, tooling

**Examples**:
```
feat: Add CSV export command

Add a new command to export data as CSV files with configurable
delimiters and encoding options.

Closes #123
```

```
fix: Handle empty input in greet command

The greet command now properly validates input and returns
a clear error message when the name is empty.
```

### Pull Request Process

1. **Fork** the repository
2. **Create a branch** from `main`
   ```bash
   git checkout -b feat/my-feature
   ```
3. **Make your changes**
   - Follow coding standards
   - Add tests
   - Update documentation
4. **Commit with clear messages**
   ```bash
   git commit -m "feat: Add my feature"
   ```
5. **Push to your fork**
   ```bash
   git push origin feat/my-feature
   ```
6. **Create Pull Request** on GitHub
   - Fill in the PR template
   - Link related issues
   - Ensure CI checks pass

### After Submitting

- **Respond to feedback**: Address review comments promptly
- **Keep PR updated**: Rebase if needed
- **Be patient**: Reviews may take time

## Testing Guidelines

### Writing Tests

All new features must include tests. Follow the table-driven test pattern:

```go
func TestMyFunction(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "result", false},
        {"empty input", "", "", true},
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

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Run with race detection
go test -race ./...

# Run specific package
go test -v ./internal/service/myservice
```

## Documentation

### Updating Documentation

When making changes that affect functionality:

1. **Update relevant docs** in `docs/` directory
2. **Update examples** if changing command syntax
3. **Update config docs** if changing configuration
4. **Update CHANGELOG.md** with your changes

### Documentation Standards

- Use clear, concise language
- Include code examples
- Follow existing documentation structure
- Add links to related documentation

## Development Workflow

### Daily Development

```bash
# 1. Make changes to code

# 2. Format code (automatic on commit)
make fmt

# 3. Run linters
make lint

# 4. Run tests
make test

# 5. Build
make build

# 6. Test binary
./build/bin/termplate <your-command>
```

### Git Hooks (Automatic)

The project uses lefthook for git hooks:

- **Pre-commit**: Automatically formats code and runs linter
- **Pre-push**: Runs tests and vet

If hooks fail, fix the issues before committing/pushing.

## Questions?

- **Check documentation**: See `docs/` directory
- **Ask in issues**: Create a question issue
- **Review examples**: Check `cmd/example/` for working code
- **Read conventions**: See [CONVENTIONS.md](CONVENTIONS.md)

## Recognition

Contributors will be recognized in the project. Thank you for your contributions! 🙏

---

**Happy contributing!** 🚀

If you're unsure about anything, just ask. We're here to help!
