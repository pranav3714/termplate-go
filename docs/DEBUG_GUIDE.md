# Debug Guide

> **Document Type**: Development Guide
> **Purpose**: Comprehensive debugging strategies and tools for Termplate Go
> **Keywords**: debug, debugging, troubleshoot, breakpoint, delve, logs, trace, inspect, gdb
> **Related**: GETTING_STARTED.md, TESTING_PATTERNS.md, PERFORMANCE_GUIDE.md

This guide covers debugging techniques, tools, and workflows for Termplate Go CLI applications.

## Quick Start

### Enable Verbose Logging

```bash
# Run with verbose flag
./build/bin/mycli -v mycommand

# Or use debug log level
./build/bin/mycli --log-level=debug mycommand

# Set via environment variable
export TERMPLATE_LOG_LEVEL=debug
./build/bin/mycli mycommand
```

### Check What's Running

```bash
# See all build artifacts
ls -la build/bin/

# Check current version
./build/bin/mycli version

# Validate configuration
./build/bin/mycli --config ~/.mycli.yaml mycommand --dry-run
```

## Debugging Tools

### 1. Delve (dlv) - Go Debugger

**Installation:**

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

**Basic Usage:**

```bash
# Debug the main package
dlv debug ./main.go -- mycommand --flag value

# Attach to running process
dlv attach $(pgrep mycli)

# Debug a test
dlv test ./internal/service/example -- -test.run TestServiceMethod
```

**Common Delve Commands:**

```
(dlv) break main.main          # Set breakpoint at main
(dlv) break handler.go:42      # Set breakpoint at line
(dlv) break Service.Execute    # Set breakpoint at method
(dlv) continue                 # Continue execution
(dlv) next                     # Step over
(dlv) step                     # Step into
(dlv) print myVar              # Print variable
(dlv) locals                   # Show local variables
(dlv) args                     # Show function arguments
(dlv) stack                    # Show stack trace
(dlv) goroutines               # List all goroutines
(dlv) quit                     # Exit debugger
```

### 2. Using Print Debugging

**Structured Logging with slog:**

```go
import "log/slog"

// In your code
func (s *Service) DoSomething(ctx context.Context, input string) error {
    slog.InfoContext(ctx, "entering DoSomething",
        "input", input,
        "timestamp", time.Now(),
    )

    // Your logic
    result := processInput(input)

    slog.DebugContext(ctx, "intermediate result",
        "result", result,
        "length", len(result),
    )

    return nil
}
```

**Output format:**

```json
{
  "time": "2026-01-18T10:30:45.123Z",
  "level": "INFO",
  "msg": "entering DoSomething",
  "input": "test-value",
  "timestamp": "2026-01-18T10:30:45.123Z"
}
```

### 3. Using fmt.Printf for Quick Debugging

```go
import "fmt"

func debugPrint(label string, value interface{}) {
    fmt.Printf("[DEBUG] %s: %+v\n", label, value)
}

// Usage
debugPrint("config", cfg)
debugPrint("result", result)
```

### 4. Go's Built-in Profiling

```bash
# CPU profiling
go build -o mycli
./mycli mycommand --cpuprofile=cpu.prof

# Memory profiling
./mycli mycommand --memprofile=mem.prof

# Analyze with pprof
go tool pprof cpu.prof
go tool pprof mem.prof
```

See **PERFORMANCE_GUIDE.md** for detailed profiling instructions.

## Debug Workflows

### Workflow 1: Command Not Working

**Steps:**

1. **Enable verbose logging**
   ```bash
   ./build/bin/mycli -v mycommand --flag value
   ```

2. **Check configuration**
   ```bash
   # See what config is loaded
   ./build/bin/mycli --config ~/.mycli.yaml -v mycommand

   # Or use debug to see all config values
   ./build/bin/mycli --log-level=debug mycommand
   ```

3. **Inspect the code path**
   ```bash
   # Add logging to trace execution
   # In cmd/mycommand/mycommand.go
   slog.Info("command started", "args", args)

   # In internal/handler/myhandler.go
   slog.Debug("handler input", "input", input)

   # In internal/service/myservice/service.go
   slog.Debug("service called", "param", param)
   ```

4. **Rebuild and test**
   ```bash
   make build
   ./build/bin/mycli -v mycommand
   ```

5. **If still failing, use delve**
   ```bash
   dlv debug ./main.go -- mycommand --flag value
   (dlv) break cmd/mycommand/mycommand.go:30
   (dlv) continue
   (dlv) print args
   (dlv) print err
   ```

### Workflow 2: Unexpected Behavior

**Steps:**

1. **Write a test to reproduce**
   ```go
   // internal/service/myservice/service_test.go
   func TestService_UnexpectedBehavior(t *testing.T) {
       s := NewService()

       got, err := s.DoSomething(context.Background(), "problematic-input")
       if err != nil {
           t.Errorf("unexpected error: %v", err)
       }

       want := "expected-output"
       if got != want {
           t.Errorf("got %v, want %v", got, want)
       }
   }
   ```

2. **Run test with verbose output**
   ```bash
   go test -v ./internal/service/myservice -run TestService_UnexpectedBehavior
   ```

3. **Debug the test**
   ```bash
   dlv test ./internal/service/myservice -- -test.run TestService_UnexpectedBehavior
   (dlv) break service.go:45
   (dlv) continue
   (dlv) print input
   (dlv) next
   (dlv) print intermediateResult
   ```

4. **Add strategic logging**
   ```go
   func (s *Service) DoSomething(ctx context.Context, input string) (string, error) {
       slog.Debug("DoSomething start", "input", input)

       // Add checkpoints
       if len(input) > 100 {
           slog.Warn("input too long", "length", len(input))
       }

       result := process(input)
       slog.Debug("after processing", "result", result)

       return result, nil
   }
   ```

### Workflow 3: Performance Issues

**Steps:**

1. **Add timing logs**
   ```go
   func (s *Service) DoSomething(ctx context.Context, input string) error {
       start := time.Now()
       defer func() {
           duration := time.Since(start)
           slog.Info("DoSomething completed",
               "duration_ms", duration.Milliseconds(),
               "input_size", len(input),
           )
       }()

       // Your logic
       return nil
   }
   ```

2. **Profile the application**
   ```bash
   # See PERFORMANCE_GUIDE.md for detailed profiling
   go test -bench=. -cpuprofile=cpu.prof ./internal/service/myservice
   go tool pprof cpu.prof
   ```

3. **Check for goroutine leaks**
   ```bash
   # Add goroutine count logging
   slog.Info("goroutine count", "count", runtime.NumGoroutine())
   ```

### Workflow 4: Panic or Crash

**Steps:**

1. **Get stack trace**
   ```bash
   # Run with GOTRACEBACK
   GOTRACEBACK=all ./build/bin/mycli mycommand
   ```

2. **Add recovery handler**
   ```go
   func safeExecute(fn func() error) (err error) {
       defer func() {
           if r := recover(); r != nil {
               err = fmt.Errorf("panic recovered: %v\nstack: %s", r, debug.Stack())
           }
       }()

       return fn()
   }
   ```

3. **Debug with delve**
   ```bash
   dlv debug ./main.go -- mycommand
   (dlv) on panic
   (dlv) continue
   # When panic occurs, examine state
   (dlv) stack
   (dlv) locals
   ```

### Workflow 5: Integration Issues

**Steps:**

1. **Enable HTTP request logging**
   ```go
   // For API clients
   import "net/http/httputil"

   func logRequest(req *http.Request) {
       dump, _ := httputil.DumpRequest(req, true)
       slog.Debug("HTTP request", "dump", string(dump))
   }

   func logResponse(resp *http.Response) {
       dump, _ := httputil.DumpResponse(resp, true)
       slog.Debug("HTTP response", "dump", string(dump))
   }
   ```

2. **Use request ID tracing**
   ```go
   func (s *Service) DoSomething(ctx context.Context, input string) error {
       requestID := generateRequestID()
       ctx = context.WithValue(ctx, "request_id", requestID)

       slog.InfoContext(ctx, "starting request", "request_id", requestID)

       // Pass ctx through all calls
       result, err := s.callExternalAPI(ctx, input)

       slog.InfoContext(ctx, "request completed", "request_id", requestID)
       return err
   }
   ```

3. **Mock external dependencies**
   ```go
   // See TESTING_PATTERNS.md for detailed mocking
   type MockAPIClient struct {
       CallFunc func(context.Context, string) (string, error)
   }

   func (m *MockAPIClient) Call(ctx context.Context, input string) (string, error) {
       if m.CallFunc != nil {
           return m.CallFunc(ctx, input)
       }
       return "", nil
   }
   ```

## Common Issues and Solutions

### Issue: "command not found"

**Problem:** Binary not in PATH or not built

**Solution:**
```bash
# Check if built
ls -la build/bin/mycli

# If not, build it
make build

# Add to PATH
export PATH=$PATH:$(pwd)/build/bin

# Or use full path
./build/bin/mycli mycommand
```

### Issue: Configuration not loading

**Problem:** Config file in wrong location or wrong format

**Solution:**
```bash
# Check config locations (in order of precedence)
ls -la ./.mycli.yaml
ls -la ~/.mycli.yaml
ls -la /etc/mycli/config.yaml

# Validate YAML syntax
yamllint ~/.mycli.yaml

# Test with explicit config
./build/bin/mycli --config=/path/to/config.yaml -v mycommand

# Check what viper is loading
# Add to cmd/root.go:
fmt.Println("Config file used:", viper.ConfigFileUsed())
fmt.Printf("All settings: %+v\n", viper.AllSettings())
```

### Issue: Tests passing but CLI fails

**Problem:** Different configuration or context

**Solution:**
```bash
# Run integration test that mimics CLI
go test -v ./cmd -run TestIntegration

# Or create integration test
# cmd/integration_test.go:
func TestCLIExecution(t *testing.T) {
    cmd := exec.Command("./build/bin/mycli", "mycommand", "--flag", "value")
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Errorf("command failed: %v\nOutput: %s", err, output)
    }
}
```

### Issue: Goroutine leaks

**Problem:** Goroutines not properly closed

**Solution:**
```bash
# Check goroutine count
go test -v ./internal/service/myservice -run TestMyService

# Add goroutine leak detection
import "go.uber.org/goleak"

func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)
}
```

### Issue: Memory leak

**Problem:** Resources not released

**Solution:**
```bash
# Profile memory usage
go test -bench=. -memprofile=mem.prof ./internal/service/myservice

# Analyze with pprof
go tool pprof -alloc_space mem.prof
(pprof) top
(pprof) list FunctionName

# Common causes:
# 1. Not closing files/connections
# 2. Holding references in maps/slices
# 3. Goroutines holding data
```

## Debugging Best Practices

### 1. Use Structured Logging

**Good:**
```go
slog.Info("user action",
    "user_id", userID,
    "action", "delete",
    "resource", resourceID,
    "duration_ms", duration.Milliseconds(),
)
```

**Bad:**
```go
fmt.Printf("User %s deleted resource %s in %v\n", userID, resourceID, duration)
```

### 2. Add Context to Errors

**Good:**
```go
if err != nil {
    return fmt.Errorf("processing user request (user_id=%s): %w", userID, err)
}
```

**Bad:**
```go
if err != nil {
    return err
}
```

### 3. Use Debug Build Tags

```go
//go:build debug

package mypackage

func init() {
    // Enable debug features only when built with -tags debug
    debugMode = true
}
```

Build with:
```bash
go build -tags debug -o mycli-debug ./main.go
```

### 4. Add Instrumentation Points

```go
type InstrumentedService struct {
    inner Service
}

func (s *InstrumentedService) DoSomething(ctx context.Context, input string) error {
    start := time.Now()

    slog.DebugContext(ctx, "service call start", "method", "DoSomething")

    err := s.inner.DoSomething(ctx, input)

    slog.DebugContext(ctx, "service call end",
        "method", "DoSomething",
        "duration_ms", time.Since(start).Milliseconds(),
        "error", err != nil,
    )

    return err
}
```

### 5. Use Table-Driven Test Cases

```go
func TestService_Debug(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
        debug   bool  // Enable debug output for this case
    }{
        {
            name:    "problematic case",
            input:   "edge-case-input",
            want:    "expected-output",
            wantErr: false,
            debug:   true,  // Enable debug for this case only
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.debug {
                slog.SetLogLoggerLevel(slog.LevelDebug)
                defer slog.SetLogLoggerLevel(slog.LevelInfo)
            }

            // Test logic
        })
    }
}
```

## IDE Integration

### VS Code

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug CLI",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/main.go",
      "args": ["mycommand", "--flag", "value"],
      "env": {
        "TERMPLATE_LOG_LEVEL": "debug"
      }
    },
    {
      "name": "Debug Test",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${workspaceFolder}/internal/service/myservice",
      "args": ["-test.run", "TestMyFunction"]
    }
  ]
}
```

### GoLand / IntelliJ

1. Right-click on `main.go` → Debug 'go build main.go'
2. Edit configuration → Add program arguments
3. Add environment variables as needed
4. Set breakpoints by clicking line numbers
5. Use Debug panel to step through code

## Advanced Techniques

### Conditional Breakpoints in Delve

```bash
dlv debug ./main.go -- mycommand

# Break only when condition is true
(dlv) break handler.go:42
(dlv) condition 1 len(input) > 100

# Continue until condition met
(dlv) continue
```

### Remote Debugging

```bash
# On remote server
dlv exec ./mycli --headless --listen=:2345 --api-version=2 -- mycommand

# On local machine
dlv connect remote-server:2345
(dlv) break main.main
(dlv) continue
```

### Core Dump Analysis

```bash
# Enable core dumps
ulimit -c unlimited

# Run program (crashes and creates core)
./build/bin/mycli mycommand

# Analyze core dump
dlv core ./build/bin/mycli core.12345
(dlv) stack
(dlv) locals
```

### Race Condition Detection

```bash
# Build with race detector
go build -race -o mycli-race ./main.go

# Run
./mycli-race mycommand

# Or in tests
go test -race ./...
```

## Troubleshooting Checklist

When debugging issues, check:

- [ ] Verbose logging enabled (`-v` flag)
- [ ] Configuration file loaded correctly
- [ ] Environment variables set appropriately
- [ ] Dependencies up to date (`go mod tidy`)
- [ ] Build successful (`make build`)
- [ ] Tests passing (`make test`)
- [ ] Linter passing (`make lint`)
- [ ] Error messages captured in logs
- [ ] Stack traces available
- [ ] Resource cleanup (files, connections closed)
- [ ] Goroutines cleaned up
- [ ] Context properly passed through calls

## Resources

- [Delve Documentation](https://github.com/go-delve/delve/tree/master/Documentation)
- [Go Diagnostics](https://go.dev/doc/diagnostics)
- [Effective Go - Debugging](https://go.dev/doc/effective_go#debugging)
- [slog Package](https://pkg.go.dev/log/slog)
- [pprof Tutorial](https://go.dev/blog/pprof)

## See Also

- **TESTING_PATTERNS.md** - Testing strategies and mocking
- **PERFORMANCE_GUIDE.md** - Profiling and optimization
- **CONVENTIONS.md** - Code standards and error handling
- **GETTING_STARTED.md** - Basic development workflow
