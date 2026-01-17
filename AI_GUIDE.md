# AI Model Guide

> **Audience**: Claude (Sonnet/Opus/Haiku), GPT, and other AI coding assistants
> **Purpose**: Optimize AI understanding and task execution in this codebase

## 🎯 Quick Start for AI Models

When you start working on this project, **follow this sequence**:

### Phase 1: Understand Context (Always do this first)

1. **Read**: `PROJECT_CONTEXT.md` - Complete project overview
   - Project type, state, architecture
   - Directory structure
   - Key files for common tasks

2. **Read**: `CONVENTIONS.md` - Coding standards
   - Architecture principles
   - Naming conventions
   - Code patterns
   - Error handling rules

3. **Read**: `QUICK_REFERENCE.md` - Fast lookups
   - Common snippets
   - File locations
   - Quick answers

### Phase 2: Task-Specific Analysis

Based on the task type, consult:

| Task Type | Documents to Read |
|-----------|------------------|
| **Adding a command** | `docs/GETTING_STARTED.md`, `cmd/example/greet.go`, `QUICK_REFERENCE.md` |
| **Configuration** | `docs/CONFIGURATION_GUIDE.md`, `internal/config/config.go` |
| **Error handling** | `internal/model/errors.go`, `CONVENTIONS.md#error-handling` |
| **Output formatting** | `internal/output/formatter.go`, `docs/CONFIGURATION_GUIDE.md` |
| **Testing** | `docs/GO_CLI_COMPREHENSIVE_REFERENCE.md#6-testing-patterns` |
| **Architecture question** | `PROJECT_CONTEXT.md`, `CONVENTIONS.md`, `docs/PROJECT_SUMMARY.md` |

### Phase 3: Examine Examples

**Always check working examples before implementing**:

- Command: `cmd/example/greet.go`
- Handler: `internal/handler/greet.go`
- Service: `internal/service/example/service.go`
- Config: `configs/config.example.yaml`

---

## 📁 Critical Files for AI Analysis

### Must-Read Files (In Order)

1. **PROJECT_CONTEXT.md** - Start here, always
2. **CONVENTIONS.md** - Coding standards
3. **internal/config/config.go** - Configuration structure
4. **internal/model/errors.go** - Error types
5. **cmd/root.go** - CLI entry point

### Example Files (Reference These)

1. **cmd/example/greet.go** - Complete command example
2. **internal/handler/greet.go** - Handler pattern
3. **internal/service/example/service.go** - Service pattern
4. **configs/config.example.yaml** - Configuration example

### Utility Files (Use These)

1. **internal/output/formatter.go** - Output formatting
2. **internal/logger/logger.go** - Logging setup
3. **pkg/version/version.go** - Version information

---

## 🧠 Architecture Understanding

### Layer Responsibilities (Strict Rules)

```
cmd/              → CLI wiring ONLY (no business logic)
                    ↓
handler/          → Validation + I/O formatting (thin layer)
                    ↓
service/          → Business logic (framework-independent, testable)
                    ↓
repository/       → Data access (abstracted)
                    ↓
model/            → Domain entities (no dependencies)
```

### Dependency Rules (Never Violate)

1. ✅ **DO**: Flow dependencies downward only
2. ✅ **DO**: Import from lower layers
3. ❌ **DON'T**: Create circular imports
4. ❌ **DON'T**: Put business logic in `cmd/`
5. ❌ **DON'T**: Import framework code in `service/`

### Import Patterns

```go
// ✅ GOOD: Handler imports service
package handler
import "github.com/blacksilver/ever-so-powerful/internal/service/myservice"

// ✅ GOOD: Service imports repository
package service
import "github.com/blacksilver/ever-so-powerful/internal/repository"

// ❌ BAD: Service imports handler (upward dependency)
package service
import "github.com/blacksilver/ever-so-powerful/internal/handler"

// ❌ BAD: Model imports service (upward dependency)
package model
import "github.com/blacksilver/ever-so-powerful/internal/service"
```

---

## 🔍 Task Analysis Workflow

### When Asked to "Add a Feature"

**Step 1: Analyze Requirements**
- What does the feature do?
- What layer(s) are involved?
- Does it need configuration?
- Does it need new domain models?

**Step 2: Read Relevant Docs**
- `docs/GETTING_STARTED.md` - How to add commands
- `CONVENTIONS.md` - Patterns to follow
- `cmd/example/` - Working example

**Step 3: Plan File Changes**
- Command: `cmd/myfeature/myfeature.go` (CREATE)
- Handler: `internal/handler/myfeature.go` (CREATE)
- Service: `internal/service/myfeature/service.go` (CREATE)
- Registration: `cmd/root.go` (MODIFY)

**Step 4: Implement**
- Follow patterns from `cmd/example/`
- Use error wrapping: `fmt.Errorf("context: %w", err)`
- Use structured logging: `slog.InfoContext(ctx, "msg", "key", val)`
- Export types properly (PascalCase for public)

**Step 5: Verify**
- Does it follow clean architecture?
- Are errors properly wrapped?
- Is logging structured?
- Are unused parameters named `_`?

---

### When Asked About "Configuration"

**Step 1: Check Current Config**
- Read: `internal/config/config.go`
- Check: `configs/config.example.yaml`

**Step 2: Determine Scope**
- Is it API configuration?
- Is it file processing?
- Is it database?
- Is it output formatting?
- Is it something new?

**Step 3: Implement**
- Add struct fields to `internal/config/config.go`
- Add defaults to `internal/config/defaults.go`
- Document in `configs/config.example.yaml`
- Update validation in `config.go` if needed

**Step 4: Reference**
- See: `docs/CONFIGURATION_GUIDE.md` for usage patterns

---

### When Asked About "Error Handling"

**Step 1: Check Domain Errors**
- Read: `internal/model/errors.go`
- Available: `ErrNotFound`, `ErrAlreadyExists`, `ErrInvalidInput`, `ErrUnauthorized`

**Step 2: Determine Error Type**
- Is it a validation error? → `model.NewValidationError("field", "message")`
- Is it an operation error? → `model.NewOperationError("op", "entity", "id", err)`
- Is it a wrapped error? → `fmt.Errorf("context: %w", err)`

**Step 3: Always Wrap**
```go
if err != nil {
    return fmt.Errorf("descriptive context: %w", err)
}
```

---

## 🤖 AI-Specific Best Practices

### When Reading Code

1. **Start with architecture**: Understand layer boundaries
2. **Check imports**: Verify no circular dependencies
3. **Follow patterns**: Match existing code style
4. **Read examples**: `cmd/example/` shows complete flow

### When Writing Code

1. **Match existing patterns**: Don't invent new structures
2. **Follow layer rules**: Respect architecture boundaries
3. **Wrap all errors**: Use `fmt.Errorf("context: %w", err)`
4. **Use structured logging**: `slog.Info("msg", "key", value)`
5. **Name unused params**: Use `_` for unused function parameters
6. **Export correctly**: Public types are PascalCase, private are camelCase

### When Suggesting Changes

1. **Preserve architecture**: Don't violate layer boundaries
2. **Maintain patterns**: Keep code consistent
3. **Check conventions**: Follow `CONVENTIONS.md`
4. **Test before commit**: Suggest running `make lint`

---

## 📊 File Location Matrix

| Task | Create Files | Modify Files | Reference Files |
|------|-------------|--------------|-----------------|
| **Add command** | `cmd/myfeature/`, `internal/handler/myhandler.go`, `internal/service/myfeature/` | `cmd/root.go` | `cmd/example/greet.go` |
| **Add config** | - | `internal/config/config.go`, `internal/config/defaults.go`, `configs/config.example.yaml` | `docs/CONFIGURATION_GUIDE.md` |
| **Add error** | - | `internal/model/errors.go` | `CONVENTIONS.md#error-handling` |
| **Add test** | `*_test.go` next to code | - | `docs/GO_CLI_COMPREHENSIVE_REFERENCE.md#6-testing-patterns` |
| **Format output** | - | - | `internal/output/formatter.go` |

---

## 🎓 Common AI Analysis Scenarios

### Scenario 1: "Add a file processing command"

**Analysis Sequence**:
1. Read `PROJECT_CONTEXT.md` - Understand architecture
2. Check `internal/config/config.go` - See FilesConfig
3. Read `cmd/example/greet.go` - Command pattern
4. Check `docs/GETTING_STARTED.md` - Step-by-step guide

**Implementation**:
- Create `cmd/process/process.go`
- Create `internal/handler/process.go`
- Create `internal/service/process/service.go`
- Use `cfg.Files.*` for configuration
- Use `formatter.Print()` for output

---

### Scenario 2: "Configure API authentication"

**Analysis Sequence**:
1. Read `internal/config/config.go` - See APIConfig
2. Check `docs/CONFIGURATION_GUIDE.md` - API section
3. Note helper: `cfg.API.GetAPIAuthHeader()`

**Usage**:
```go
cfg, _ := config.Load()
headerName, headerValue := cfg.API.GetAPIAuthHeader()
req.Header.Set(headerName, headerValue)
```

---

### Scenario 3: "Add database operations"

**Analysis Sequence**:
1. Read `internal/config/config.go` - See DBConfig
2. Check helper: `cfg.Database.GetDSN()`
3. Create repository interface
4. Implement in `internal/repository/myrepo/`

**Pattern**:
```go
// Define interface
type Repository interface {
    Get(ctx context.Context, id string) (*model.Entity, error)
}

// Use in service
type Service struct {
    repo Repository
}
```

---

### Scenario 4: "Format output as table/JSON/CSV"

**Analysis Sequence**:
1. Check `internal/output/formatter.go` - Formatter implementation
2. Read `docs/CONFIGURATION_GUIDE.md` - Output section
3. Use `formatter.Print(data)` - Automatic formatting

**Supported data types**:
- `map[string]string`
- `[]map[string]string`
- `[][]string`

---

## ⚠️ Common Pitfalls for AI Models

### Pitfall 1: Putting Business Logic in cmd/

```go
// ❌ BAD: Business logic in command
var Cmd = &cobra.Command{
    RunE: func(cmd *cobra.Command, _ []string) error {
        // Business logic here - DON'T DO THIS
        result := processData(input)
        return nil
    },
}

// ✅ GOOD: Delegate to handler
var Cmd = &cobra.Command{
    RunE: func(cmd *cobra.Command, _ []string) error {
        h := handler.NewMyHandler()
        return h.Execute(cmd.Context(), handler.MyInput{})
    },
}
```

### Pitfall 2: Not Wrapping Errors

```go
// ❌ BAD: No context
if err != nil {
    return err
}

// ✅ GOOD: Wrapped with context
if err != nil {
    return fmt.Errorf("processing file: %w", err)
}
```

### Pitfall 3: Circular Dependencies

```go
// ❌ BAD: Handler imports service, service imports handler
// This creates a circular dependency

// ✅ GOOD: Dependencies only flow downward
// cmd → handler → service → repository → model
```

### Pitfall 4: Ignoring Existing Patterns

```go
// ❌ BAD: Creating new logging approach
fmt.Printf("Info: %s\n", message)

// ✅ GOOD: Use existing structured logging
slog.Info("operation", "detail", message)
```

---

## 📖 Documentation Discovery Map

```
Start Here → PROJECT_CONTEXT.md
                ↓
        What kind of task?
                ↓
    ┌───────────┼───────────┐
    │           │           │
    ↓           ↓           ↓
Add Feature  Configure  Understand
    │           │           │
    ↓           ↓           ↓
GETTING_      CONFIG_    CONVENTIONS.md
STARTED.md    GUIDE.md   PROJECT_SUMMARY.md
    │           │           │
    ↓           ↓           ↓
QUICK_        configs/    GO_CLI_
REFERENCE.md  example     REFERENCE.md
              .yaml
```

---

## 🔑 Key Points for AI Success

1. **Always read PROJECT_CONTEXT.md first** - Gives you full context
2. **Follow existing patterns** - Check `cmd/example/` before implementing
3. **Respect layer boundaries** - No upward dependencies
4. **Wrap all errors** - Use `fmt.Errorf("context: %w", err)`
5. **Use structured logging** - `slog.Info()` with key-value pairs
6. **Check conventions** - Before suggesting code changes
7. **Test your understanding** - Explain your approach before implementing
8. **Reference docs** - Link to relevant docs in explanations

---

## 🎯 Analysis Checklist

Before implementing a feature:

- [ ] Read `PROJECT_CONTEXT.md`
- [ ] Understand which layers are involved
- [ ] Check existing examples in `cmd/example/`
- [ ] Review relevant documentation
- [ ] Verify no circular dependencies
- [ ] Follow naming conventions
- [ ] Plan test strategy

During implementation:

- [ ] Follow clean architecture
- [ ] Wrap all errors
- [ ] Use structured logging
- [ ] Export types correctly
- [ ] Name unused params with `_`
- [ ] Add comments for exported items

After implementation:

- [ ] Verify architecture compliance
- [ ] Check error wrapping
- [ ] Ensure tests follow table-driven pattern
- [ ] Suggest running `make lint`
- [ ] Reference documentation

---

## 💡 Tips for Effective AI Assistance

### Understanding User Intent

When a user asks to:

- **"Add a feature"** → They want: cmd + handler + service + tests
- **"Configure X"** → They want: Update config.go + defaults.go + example.yaml
- **"Fix error handling"** → Check: Error wrapping, domain errors
- **"Add output format"** → Use: `internal/output/formatter.go`
- **"Test the code"** → Follow: Table-driven test pattern

### Suggesting Code

- **Always provide context**: "This follows the pattern in cmd/example/greet.go"
- **Reference docs**: "As shown in CONVENTIONS.md..."
- **Explain architecture**: "This handler delegates to the service layer because..."
- **Link to examples**: "Similar to the implementation in..."

### Explaining Architecture

- **Use diagrams**: The ASCII diagrams in PROJECT_CONTEXT.md
- **Show flow**: cmd → handler → service → repository
- **Explain why**: "This separation allows testing business logic without CLI framework"

---

## 📝 Documentation Keywords for Search

**When searching for information, these keywords map to documents:**

| Keywords | Document |
|----------|----------|
| architecture, layers, clean architecture | PROJECT_CONTEXT.md, CONVENTIONS.md |
| add command, create command, new command | docs/GETTING_STARTED.md, QUICK_REFERENCE.md |
| config, configuration, settings | docs/CONFIGURATION_GUIDE.md, internal/config/config.go |
| error, error handling, wrap error | CONVENTIONS.md#error-handling, internal/model/errors.go |
| log, logging, slog | CONVENTIONS.md#logging, internal/logger/logger.go |
| output, format, json, yaml, table, csv | internal/output/formatter.go, docs/CONFIGURATION_GUIDE.md |
| test, testing, table-driven | docs/GO_CLI_COMPREHENSIVE_REFERENCE.md#6-testing-patterns |
| cobra, command, flag | docs/GO_CLI_COMPREHENSIVE_REFERENCE.md#2-cli-framework-cobra |
| viper, configuration | docs/GO_CLI_COMPREHENSIVE_REFERENCE.md#3-configuration-viper |

---

## 🚀 Example AI Workflow

### Task: "Add a command to process CSV files"

**Step 1: Understand Context**
```
Read: PROJECT_CONTEXT.md
→ Understand: cmd → handler → service architecture
→ Note: FilesConfig already exists in internal/config/config.go
```

**Step 2: Check Examples**
```
Read: cmd/example/greet.go
→ See command structure
→ Note pattern: flags → handler → service
```

**Step 3: Plan Implementation**
```
Create:
  1. cmd/process/process.go (command definition)
  2. internal/handler/process.go (validation, I/O)
  3. internal/service/process/service.go (CSV parsing logic)

Modify:
  1. cmd/root.go (register command)
```

**Step 4: Implement**
```go
// Follow patterns from cmd/example/greet.go
// Use cfg.Files.* for file configuration
// Use formatter.Print() for output
// Wrap errors: fmt.Errorf("context: %w", err)
// Log: slog.InfoContext(ctx, "processing", "file", filename)
```

**Step 5: Suggest Tests**
```
Create: internal/service/process/service_test.go
Follow: Table-driven test pattern from docs/GO_CLI_COMPREHENSIVE_REFERENCE.md
```

---

## 🎯 AI Success Metrics

A successful AI interaction with this codebase:

✅ References relevant documentation
✅ Follows existing patterns (cmd/example/)
✅ Respects architecture boundaries
✅ Wraps all errors properly
✅ Uses structured logging
✅ Suggests appropriate tests
✅ Explains reasoning with doc references

---

## 📚 Quick Reference Links

| Need | File |
|------|------|
| Project overview | `PROJECT_CONTEXT.md` |
| Coding standards | `CONVENTIONS.md` |
| Quick snippets | `QUICK_REFERENCE.md` |
| Add features | `docs/GETTING_STARTED.md` |
| Configuration | `docs/CONFIGURATION_GUIDE.md` |
| CLI patterns | `docs/GO_CLI_COMPREHENSIVE_REFERENCE.md` |
| Working example | `cmd/example/greet.go` |
| All docs | `docs/README.md` |

---

## 🌟 Final Note for AI Models

This project is **well-documented** and **follows strict conventions**.

**Your job is to**:
1. Understand the architecture (PROJECT_CONTEXT.md)
2. Follow the patterns (CONVENTIONS.md)
3. Reference the docs (docs/)
4. Maintain consistency

**Success means**: Code that looks like it was written by the same person who wrote the existing code.

**Read the docs first. They will save you time and ensure correct implementation.**

---

**Happy analyzing! 🤖**
