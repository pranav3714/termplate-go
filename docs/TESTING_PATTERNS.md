# Testing Patterns Guide

> **Document Type**: Development Guide
> **Purpose**: Comprehensive testing strategies, mocking patterns, and best practices
> **Keywords**: test, testing, mock, stub, integration, unit, table-driven, coverage, tdd
> **Related**: DEBUG_GUIDE.md, GETTING_STARTED.md, CONVENTIONS.md

This guide covers testing strategies, mocking patterns, and best practices for Termplate Go CLI applications.

## Quick Start

### Run Tests

```bash
# All tests
make test

# With coverage
make coverage

# Specific package
go test -v ./internal/service/myservice

# With race detection
go test -race ./...

# Run only specific test
go test -v ./internal/service/myservice -run TestServiceMethod

# Run tests matching pattern
go test -v ./... -run TestService

# Show test output even on success
go test -v ./...

# Parallel execution (default)
go test -parallel 4 ./...
```

### Coverage Report

```bash
# Generate coverage
make coverage

# View HTML report
go tool cover -html=coverage.out

# Show coverage by function
go tool cover -func=coverage.out

# Require minimum coverage
go test -cover -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//'
```

## Testing Philosophy

### Testing Pyramid

```
       /\
      /  \      E2E Tests (Few)
     /____\     - Full CLI execution
    /      \    - Integration with real systems
   /        \
  /__________\  Integration Tests (Some)
 /            \ - Multiple components
/              \- Mock external dependencies
/______________\
                Unit Tests (Many)
                - Single functions/methods
                - Fast, isolated
                - High coverage
```

### What to Test

**DO Test:**
- Business logic in services
- Input validation in handlers
- Error conditions and edge cases
- Configuration parsing
- Data transformations
- Complex algorithms

**DON'T Test:**
- Cobra command setup (framework code)
- Simple getters/setters
- Generated code
- Third-party libraries

## Unit Testing Patterns

### 1. Table-Driven Tests

**Pattern:**

```go
func TestService_ProcessInput(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   "hello",
            want:    "HELLO",
            wantErr: false,
        },
        {
            name:    "empty input",
            input:   "",
            want:    "",
            wantErr: true,
        },
        {
            name:    "special characters",
            input:   "hello@world!",
            want:    "HELLO@WORLD!",
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := NewService()
            got, err := s.ProcessInput(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessInput() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if got != tt.want {
                t.Errorf("ProcessInput() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Benefits:**
- Easy to add new test cases
- Clear test intent
- Parallel execution friendly
- Comprehensive coverage

### 2. Parallel Test Execution

**Pattern:**

```go
func TestService_Parallel(t *testing.T) {
    tests := []struct {
        name  string
        input string
    }{
        {"case1", "input1"},
        {"case2", "input2"},
        {"case3", "input3"},
    }

    for _, tt := range tests {
        tt := tt // Capture range variable
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel() // Enable parallel execution

            s := NewService()
            result, err := s.Process(tt.input)

            if err != nil {
                t.Errorf("unexpected error: %v", err)
            }

            // Assertions
        })
    }
}
```

**Important:** Always capture loop variables when using `t.Parallel()`

### 3. Testing with Context

**Pattern:**

```go
func TestService_WithContext(t *testing.T) {
    tests := []struct {
        name        string
        ctx         context.Context
        input       string
        wantErr     bool
        wantTimeout bool
    }{
        {
            name:    "normal context",
            ctx:     context.Background(),
            input:   "test",
            wantErr: false,
        },
        {
            name: "cancelled context",
            ctx: func() context.Context {
                ctx, cancel := context.WithCancel(context.Background())
                cancel()
                return ctx
            }(),
            input:   "test",
            wantErr: true,
        },
        {
            name: "timeout context",
            ctx: func() context.Context {
                ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
                defer cancel()
                time.Sleep(10 * time.Millisecond)
                return ctx
            }(),
            input:       "test",
            wantErr:     true,
            wantTimeout: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := NewService()
            _, err := s.ProcessWithContext(tt.ctx, tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }

            if tt.wantTimeout && !errors.Is(err, context.DeadlineExceeded) {
                t.Errorf("expected timeout error, got %v", err)
            }
        })
    }
}
```

### 4. Testing Error Types

**Pattern:**

```go
func TestService_ErrorTypes(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        wantError error
    }{
        {
            name:      "validation error",
            input:     "",
            wantError: model.ErrInvalidInput,
        },
        {
            name:      "not found error",
            input:     "nonexistent",
            wantError: model.ErrNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := NewService()
            _, err := s.Process(tt.input)

            if !errors.Is(err, tt.wantError) {
                t.Errorf("error = %v, want %v", err, tt.wantError)
            }
        })
    }
}
```

## Mocking Patterns

### 1. Interface-Based Mocking

**Define interface:**

```go
// internal/service/myservice/repository.go
package myservice

import "context"

type Repository interface {
    Get(ctx context.Context, id string) (*Model, error)
    Save(ctx context.Context, model *Model) error
    Delete(ctx context.Context, id string) error
}
```

**Implement service with interface:**

```go
// internal/service/myservice/service.go
package myservice

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) Process(ctx context.Context, id string) error {
    model, err := s.repo.Get(ctx, id)
    if err != nil {
        return err
    }

    // Process model
    model.Status = "processed"

    return s.repo.Save(ctx, model)
}
```

**Create mock in test:**

```go
// internal/service/myservice/service_test.go
package myservice

import (
    "context"
    "errors"
    "testing"
)

type MockRepository struct {
    GetFunc    func(ctx context.Context, id string) (*Model, error)
    SaveFunc   func(ctx context.Context, model *Model) error
    DeleteFunc func(ctx context.Context, id string) error
}

func (m *MockRepository) Get(ctx context.Context, id string) (*Model, error) {
    if m.GetFunc != nil {
        return m.GetFunc(ctx, id)
    }
    return nil, errors.New("not implemented")
}

func (m *MockRepository) Save(ctx context.Context, model *Model) error {
    if m.SaveFunc != nil {
        return m.SaveFunc(ctx, model)
    }
    return errors.New("not implemented")
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
    if m.DeleteFunc != nil {
        return m.DeleteFunc(ctx, id)
    }
    return errors.New("not implemented")
}

func TestService_Process(t *testing.T) {
    mock := &MockRepository{
        GetFunc: func(ctx context.Context, id string) (*Model, error) {
            return &Model{ID: id, Status: "pending"}, nil
        },
        SaveFunc: func(ctx context.Context, model *Model) error {
            if model.Status != "processed" {
                t.Errorf("expected status 'processed', got '%s'", model.Status)
            }
            return nil
        },
    }

    service := NewService(mock)
    err := service.Process(context.Background(), "test-id")

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}
```

### 2. Spy Pattern (Tracking Calls)

**Pattern:**

```go
type SpyRepository struct {
    GetCalls    []string
    SaveCalls   []*Model
    DeleteCalls []string
}

func (s *SpyRepository) Get(ctx context.Context, id string) (*Model, error) {
    s.GetCalls = append(s.GetCalls, id)
    return &Model{ID: id}, nil
}

func (s *SpyRepository) Save(ctx context.Context, model *Model) error {
    s.SaveCalls = append(s.SaveCalls, model)
    return nil
}

func (s *SpyRepository) Delete(ctx context.Context, id string) error {
    s.DeleteCalls = append(s.DeleteCalls, id)
    return nil
}

func TestService_CallTracking(t *testing.T) {
    spy := &SpyRepository{}
    service := NewService(spy)

    _ = service.Process(context.Background(), "id1")
    _ = service.Process(context.Background(), "id2")

    if len(spy.GetCalls) != 2 {
        t.Errorf("expected 2 Get calls, got %d", len(spy.GetCalls))
    }

    if len(spy.SaveCalls) != 2 {
        t.Errorf("expected 2 Save calls, got %d", len(spy.SaveCalls))
    }
}
```

### 3. HTTP Client Mocking

**Pattern:**

```go
// internal/client/client.go
package client

import (
    "context"
    "net/http"
)

type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}

type Client struct {
    http HTTPClient
}

func NewClient(httpClient HTTPClient) *Client {
    if httpClient == nil {
        httpClient = &http.Client{}
    }
    return &Client{http: httpClient}
}

// client_test.go
type MockHTTPClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    if m.DoFunc != nil {
        return m.DoFunc(req)
    }
    return nil, errors.New("not implemented")
}

func TestClient_Request(t *testing.T) {
    mock := &MockHTTPClient{
        DoFunc: func(req *http.Request) (*http.Response, error) {
            return &http.Response{
                StatusCode: 200,
                Body:       io.NopCloser(strings.NewReader(`{"status":"ok"}`)),
            }, nil
        },
    }

    client := NewClient(mock)
    // Test client methods
}
```

### 4. Time Mocking

**Pattern:**

```go
// internal/service/myservice/time.go
package myservice

import "time"

// Clock interface for mocking time
type Clock interface {
    Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time {
    return time.Now()
}

type Service struct {
    clock Clock
}

func NewService() *Service {
    return &Service{clock: realClock{}}
}

// service_test.go
type MockClock struct {
    NowFunc func() time.Time
}

func (m MockClock) Now() time.Time {
    if m.NowFunc != nil {
        return m.NowFunc()
    }
    return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
}

func TestService_WithTime(t *testing.T) {
    fixedTime := time.Date(2026, 1, 18, 12, 0, 0, 0, time.UTC)

    mock := MockClock{
        NowFunc: func() time.Time {
            return fixedTime
        },
    }

    service := &Service{clock: mock}

    result := service.DoSomethingWithTime()

    // Assert result uses fixedTime
}
```

## Integration Testing

### 1. Command Integration Tests

**Pattern:**

```go
// cmd/mycommand/integration_test.go
package mycommand_test

import (
    "bytes"
    "context"
    "testing"

    "github.com/yourorg/yourproject/cmd/mycommand"
    "github.com/spf13/cobra"
)

func TestCommand_Integration(t *testing.T) {
    tests := []struct {
        name       string
        args       []string
        wantOutput string
        wantErr    bool
    }{
        {
            name:       "valid command",
            args:       []string{"--name", "test"},
            wantOutput: "Success",
            wantErr:    false,
        },
        {
            name:    "missing required flag",
            args:    []string{},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := mycommand.Cmd
            cmd.SetArgs(tt.args)

            var output bytes.Buffer
            cmd.SetOut(&output)
            cmd.SetErr(&output)

            err := cmd.ExecuteContext(context.Background())

            if (err != nil) != tt.wantErr {
                t.Errorf("command error = %v, wantErr %v", err, tt.wantErr)
            }

            if tt.wantOutput != "" && !strings.Contains(output.String(), tt.wantOutput) {
                t.Errorf("output = %v, want contains %v", output.String(), tt.wantOutput)
            }
        })
    }
}
```

### 2. End-to-End CLI Tests

**Pattern:**

```go
// test/e2e/cli_test.go
//go:build e2e

package e2e_test

import (
    "os/exec"
    "strings"
    "testing"
)

func TestCLI_E2E(t *testing.T) {
    tests := []struct {
        name       string
        args       []string
        wantOutput string
        wantErr    bool
    }{
        {
            name:       "version command",
            args:       []string{"version"},
            wantOutput: "v0.1.0",
            wantErr:    false,
        },
        {
            name:       "help command",
            args:       []string{"--help"},
            wantOutput: "Usage:",
            wantErr:    false,
        },
    }

    // Build binary first
    buildCmd := exec.Command("make", "build")
    if err := buildCmd.Run(); err != nil {
        t.Fatalf("failed to build binary: %v", err)
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := exec.Command("./build/bin/mycli", tt.args...)
            output, err := cmd.CombinedOutput()

            if (err != nil) != tt.wantErr {
                t.Errorf("command error = %v, wantErr %v\nOutput: %s", err, tt.wantErr, output)
            }

            if !strings.Contains(string(output), tt.wantOutput) {
                t.Errorf("output = %s, want contains %s", output, tt.wantOutput)
            }
        })
    }
}
```

**Run E2E tests:**

```bash
go test -v -tags=e2e ./test/e2e
```

### 3. Database Integration Tests

**Pattern:**

```go
//go:build integration

package repository_test

import (
    "context"
    "database/sql"
    "testing"

    _ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("postgres", "postgres://localhost/test?sslmode=disable")
    if err != nil {
        t.Fatalf("failed to connect to test db: %v", err)
    }

    // Run migrations
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS items (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL
    )`)
    if err != nil {
        t.Fatalf("failed to create table: %v", err)
    }

    t.Cleanup(func() {
        db.Exec("DROP TABLE items")
        db.Close()
    })

    return db
}

func TestRepository_Integration(t *testing.T) {
    db := setupTestDB(t)
    repo := NewRepository(db)

    // Test create
    err := repo.Save(context.Background(), &Item{ID: "1", Name: "test"})
    if err != nil {
        t.Errorf("Save failed: %v", err)
    }

    // Test get
    item, err := repo.Get(context.Background(), "1")
    if err != nil {
        t.Errorf("Get failed: %v", err)
    }

    if item.Name != "test" {
        t.Errorf("got name %s, want test", item.Name)
    }
}
```

**Run integration tests:**

```bash
go test -v -tags=integration ./...
```

## Test Helpers

### 1. Setup and Teardown

**Pattern:**

```go
func TestMain(m *testing.M) {
    // Global setup
    setup()

    // Run tests
    code := m.Run()

    // Global teardown
    teardown()

    os.Exit(code)
}

func setup() {
    // Initialize test resources
    os.Setenv("ENV", "test")
}

func teardown() {
    // Cleanup test resources
    os.Unsetenv("ENV")
}

func TestWithCleanup(t *testing.T) {
    // Test-specific setup
    tempFile, err := os.CreateTemp("", "test-*")
    if err != nil {
        t.Fatal(err)
    }

    // Register cleanup
    t.Cleanup(func() {
        os.Remove(tempFile.Name())
    })

    // Test logic
}
```

### 2. Test Fixtures

**Pattern:**

```go
// testdata/fixtures.go
package testdata

func ValidInput() *Input {
    return &Input{
        Name:  "test",
        Value: 123,
    }
}

func InvalidInput() *Input {
    return &Input{
        Name:  "",
        Value: -1,
    }
}

// Usage in tests
func TestService_WithFixtures(t *testing.T) {
    input := testdata.ValidInput()

    service := NewService()
    result, err := service.Process(input)

    // Assertions
}
```

### 3. Golden Files

**Pattern:**

```go
func TestOutput_Golden(t *testing.T) {
    tests := []struct {
        name  string
        input string
    }{
        {"case1", "input1"},
        {"case2", "input2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := generateOutput(tt.input)

            goldenFile := filepath.Join("testdata", tt.name+".golden")

            if *update {
                // Update golden file
                os.WriteFile(goldenFile, []byte(got), 0644)
                return
            }

            // Compare with golden file
            want, err := os.ReadFile(goldenFile)
            if err != nil {
                t.Fatal(err)
            }

            if got != string(want) {
                t.Errorf("got:\n%s\n\nwant:\n%s", got, want)
            }
        })
    }
}

var update = flag.Bool("update", false, "update golden files")
```

**Usage:**

```bash
# Run tests normally
go test ./...

# Update golden files
go test ./... -update
```

## Test Organization

### Directory Structure

```
myproject/
├── internal/
│   └── service/
│       └── myservice/
│           ├── service.go
│           ├── service_test.go       # Unit tests
│           └── testdata/
│               ├── fixtures.go       # Test data
│               └── golden/           # Golden files
├── test/
│   ├── integration/                  # Integration tests
│   │   └── service_integration_test.go
│   ├── e2e/                          # E2E tests
│   │   └── cli_test.go
│   └── testutil/                     # Test utilities
│       ├── mocks.go                  # Shared mocks
│       └── helpers.go                # Test helpers
```

### Test File Naming

- Unit tests: `*_test.go` (same package)
- Integration tests: `*_integration_test.go` (build tag: `integration`)
- E2E tests: `*_e2e_test.go` (build tag: `e2e`)
- Test helpers: `testutil`, `testdata` directories

## Best Practices

### 1. Test Independence

**Good:**
```go
func TestA(t *testing.T) {
    service := NewService()
    // Test logic
}

func TestB(t *testing.T) {
    service := NewService()
    // Test logic - independent of TestA
}
```

**Bad:**
```go
var sharedService *Service

func TestA(t *testing.T) {
    sharedService = NewService()
    // Test logic
}

func TestB(t *testing.T) {
    // Uses sharedService - depends on TestA
}
```

### 2. Clear Test Names

**Good:**
```go
func TestService_ProcessInput_ReturnsErrorForEmptyInput(t *testing.T)
func TestService_ProcessInput_ConvertsToUpperCase(t *testing.T)
func TestService_ProcessInput_TrimsWhitespace(t *testing.T)
```

**Bad:**
```go
func TestProcess(t *testing.T)
func TestService1(t *testing.T)
func TestEdgeCase(t *testing.T)
```

### 3. Assertion Messages

**Good:**
```go
if got != want {
    t.Errorf("ProcessInput(%q) = %q, want %q", input, got, want)
}
```

**Bad:**
```go
if got != want {
    t.Errorf("test failed")
}
```

### 4. Use Subtests

**Good:**
```go
func TestService(t *testing.T) {
    t.Run("valid input", func(t *testing.T) { /* ... */ })
    t.Run("invalid input", func(t *testing.T) { /* ... */ })
    t.Run("edge case", func(t *testing.T) { /* ... */ })
}
```

### 5. Test Error Paths

```go
func TestService_ErrorHandling(t *testing.T) {
    tests := []struct {
        name      string
        setupMock func(*MockRepo)
        wantErr   error
    }{
        {
            name: "repository error",
            setupMock: func(m *MockRepo) {
                m.GetFunc = func(ctx context.Context, id string) (*Model, error) {
                    return nil, errors.New("db error")
                }
            },
            wantErr: ErrDatabase,
        },
    }
    // Test logic
}
```

## Coverage Goals

- **Overall**: 70-80% coverage
- **Business Logic (services)**: 90%+ coverage
- **Handlers**: 80%+ coverage
- **Commands**: 50-60% coverage (framework code)

**Check coverage:**

```bash
make coverage
go tool cover -func=coverage.out | grep total
```

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table-Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Testing in Go](https://go.dev/doc/tutorial/add-a-test)
- [Go Test Comments](https://github.com/golang/go/wiki/TestComments)

## See Also

- **DEBUG_GUIDE.md** - Debugging strategies
- **CONVENTIONS.md** - Code standards
- **GETTING_STARTED.md** - Development workflow
