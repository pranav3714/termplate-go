# Performance Guide

> **Document Type**: Development Guide
> **Purpose**: Profiling, benchmarking, and optimization strategies for Termplate Go
> **Keywords**: performance, profiling, benchmark, optimization, pprof, memory, cpu, trace
> **Related**: DEBUG_GUIDE.md, TESTING_PATTERNS.md, DOCKER_DEVELOPMENT.md

This guide covers performance profiling, benchmarking, and optimization techniques for Termplate Go CLI applications.

## Quick Start

### CPU Profiling

```bash
# Build with optimizations
go build -o mycli ./main.go

# Run with CPU profiling
./mycli mycommand --cpuprofile=cpu.prof

# Analyze profile
go tool pprof cpu.prof

# Interactive commands:
(pprof) top        # Show top functions by CPU time
(pprof) list main  # Show source code with annotations
(pprof) web        # Generate visual graph (requires graphviz)
```

### Memory Profiling

```bash
# Run with memory profiling
./mycli mycommand --memprofile=mem.prof

# Analyze memory allocations
go tool pprof -alloc_space mem.prof
(pprof) top

# Analyze memory in use
go tool pprof -inuse_space mem.prof
(pprof) top
```

### Quick Benchmark

```bash
# Run benchmarks
go test -bench=. -benchmem ./internal/service/myservice

# Run specific benchmark
go test -bench=BenchmarkProcess -benchmem ./internal/service/myservice

# With CPU profile
go test -bench=. -cpuprofile=cpu.prof ./internal/service/myservice
go tool pprof cpu.prof
```

## Profiling Tools

### 1. CPU Profiling with pprof

**Add profiling to your command:**

```go
// cmd/mycommand/mycommand.go
package mycommand

import (
    "os"
    "runtime/pprof"

    "github.com/spf13/cobra"
)

var cpuProfile string

var Cmd = &cobra.Command{
    Use: "mycommand",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Start CPU profiling
        if cpuProfile != "" {
            f, err := os.Create(cpuProfile)
            if err != nil {
                return err
            }
            defer f.Close()

            if err := pprof.StartCPUProfile(f); err != nil {
                return err
            }
            defer pprof.StopCPUProfile()
        }

        // Your command logic
        return nil
    },
}

func init() {
    Cmd.Flags().StringVar(&cpuProfile, "cpuprofile", "", "write CPU profile to file")
}
```

**Usage:**

```bash
# Run with profiling
./mycli mycommand --cpuprofile=cpu.prof

# Analyze
go tool pprof cpu.prof

# Common pprof commands
(pprof) top           # Top functions by CPU time
(pprof) top -cum      # Top functions by cumulative time
(pprof) list FuncName # Source code with annotations
(pprof) web           # Visual graph (requires graphviz)
(pprof) pdf > cpu.pdf # Export to PDF
```

### 2. Memory Profiling

**Add memory profiling:**

```go
// cmd/mycommand/mycommand.go
var memProfile string

var Cmd = &cobra.Command{
    Use: "mycommand",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Your command logic
        result := doWork()

        // Write memory profile
        if memProfile != "" {
            f, err := os.Create(memProfile)
            if err != nil {
                return err
            }
            defer f.Close()

            runtime.GC() // Get up-to-date statistics
            if err := pprof.WriteHeapProfile(f); err != nil {
                return err
            }
        }

        return nil
    },
}

func init() {
    Cmd.Flags().StringVar(&memProfile, "memprofile", "", "write memory profile to file")
}
```

**Analyze memory:**

```bash
# Allocations (what was allocated)
go tool pprof -alloc_space mem.prof

# In-use memory (what's still in use)
go tool pprof -inuse_space mem.prof

# Common commands
(pprof) top
(pprof) list FuncName
(pprof) web
```

### 3. Execution Tracing

**Add trace support:**

```go
import (
    "os"
    "runtime/trace"
)

var traceFile string

var Cmd = &cobra.Command{
    Use: "mycommand",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Start tracing
        if traceFile != "" {
            f, err := os.Create(traceFile)
            if err != nil {
                return err
            }
            defer f.Close()

            if err := trace.Start(f); err != nil {
                return err
            }
            defer trace.Stop()
        }

        // Your logic
        return nil
    },
}

func init() {
    Cmd.Flags().StringVar(&traceFile, "trace", "", "write execution trace to file")
}
```

**Analyze trace:**

```bash
# Run with tracing
./mycli mycommand --trace=trace.out

# View trace in browser
go tool trace trace.out

# Shows:
# - Goroutine activity
# - GC events
# - Network blocking
# - Syscalls
# - Goroutine analysis
```

### 4. HTTP pprof Server

**For long-running services:**

```go
import (
    "net/http"
    _ "net/http/pprof"
)

func main() {
    // Start pprof server
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()

    // Your application logic
}
```

**Access profiles:**

```bash
# CPU profile (30 seconds)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine profile
go tool pprof http://localhost:6060/debug/pprof/goroutine

# View in browser
open http://localhost:6060/debug/pprof
```

## Benchmarking

### Writing Benchmarks

**Basic benchmark:**

```go
// internal/service/myservice/service_bench_test.go
package myservice

import "testing"

func BenchmarkService_Process(b *testing.B) {
    s := NewService()
    input := "test-input"

    b.ResetTimer() // Reset timer after setup

    for i := 0; i < b.N; i++ {
        _, err := s.Process(context.Background(), input)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

**Run benchmarks:**

```bash
# Run all benchmarks
go test -bench=. ./internal/service/myservice

# With memory statistics
go test -bench=. -benchmem ./internal/service/myservice

# Run specific benchmark
go test -bench=BenchmarkProcess ./internal/service/myservice

# Multiple runs for accuracy
go test -bench=. -benchtime=10s ./internal/service/myservice
go test -bench=. -count=5 ./internal/service/myservice
```

### Table-Driven Benchmarks

```go
func BenchmarkProcess(b *testing.B) {
    tests := []struct {
        name  string
        input string
    }{
        {"small", "test"},
        {"medium", strings.Repeat("test", 100)},
        {"large", strings.Repeat("test", 10000)},
    }

    for _, tt := range tests {
        b.Run(tt.name, func(b *testing.B) {
            s := NewService()

            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                s.Process(context.Background(), tt.input)
            }
        })
    }
}
```

**Output:**

```
BenchmarkProcess/small-8     1000000    1234 ns/op    128 B/op    2 allocs/op
BenchmarkProcess/medium-8     100000   12345 ns/op   1024 B/op   10 allocs/op
BenchmarkProcess/large-8       10000  123456 ns/op  10240 B/op  100 allocs/op
```

### Parallel Benchmarks

```go
func BenchmarkProcessParallel(b *testing.B) {
    s := NewService()
    input := "test-input"

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            s.Process(context.Background(), input)
        }
    })
}
```

### Benchmark Comparison

**Save baseline:**

```bash
go test -bench=. -benchmem ./internal/service/myservice > old.txt
```

**Make changes, then compare:**

```bash
go test -bench=. -benchmem ./internal/service/myservice > new.txt

# Install benchstat
go install golang.org/x/perf/cmd/benchstat@latest

# Compare
benchstat old.txt new.txt
```

**Output:**

```
name        old time/op  new time/op  delta
Process-8   1.23µs ± 2%  0.98µs ± 1%  -20.33%  (p=0.000 n=10+10)

name        old alloc/op  new alloc/op  delta
Process-8   128B ± 0%      96B ± 0%     -25.00%  (p=0.000 n=10+10)
```

## Optimization Techniques

### 1. Reduce Allocations

**Problem: Too many allocations**

```go
// Bad: Creates new slice on each append
func processItems(items []string) []string {
    var result []string
    for _, item := range items {
        result = append(result, process(item))
    }
    return result
}
```

**Solution: Pre-allocate capacity**

```go
// Good: Pre-allocate with known capacity
func processItems(items []string) []string {
    result := make([]string, 0, len(items))
    for _, item := range items {
        result = append(result, process(item))
    }
    return result
}
```

### 2. Reuse Buffers

**Problem: Creating temporary buffers**

```go
// Bad: Creates new buffer each time
func formatOutput(data string) string {
    var buf bytes.Buffer
    buf.WriteString("[")
    buf.WriteString(data)
    buf.WriteString("]")
    return buf.String()
}
```

**Solution: Use sync.Pool**

```go
// Good: Reuse buffers
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func formatOutput(data string) string {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()

    buf.WriteString("[")
    buf.WriteString(data)
    buf.WriteString("]")
    return buf.String()
}
```

### 3. Avoid String Concatenation

**Problem: String concatenation creates copies**

```go
// Bad: Creates multiple string copies
func buildMessage(parts []string) string {
    msg := ""
    for _, part := range parts {
        msg += part + " "
    }
    return msg
}
```

**Solution: Use strings.Builder**

```go
// Good: Single allocation
func buildMessage(parts []string) string {
    var sb strings.Builder
    sb.Grow(len(parts) * 10) // Estimate size

    for _, part := range parts {
        sb.WriteString(part)
        sb.WriteString(" ")
    }
    return sb.String()
}
```

### 4. Efficient Map Operations

**Problem: Map lookups in loops**

```go
// Bad: Multiple lookups
func updateMap(m map[string]int, key string) {
    if _, ok := m[key]; ok {
        m[key]++
    } else {
        m[key] = 1
    }
}
```

**Solution: Single lookup**

```go
// Good: Single lookup
func updateMap(m map[string]int, key string) {
    m[key]++
}
```

### 5. Use Pointers for Large Structs

**Problem: Copying large structs**

```go
// Bad: Passes by value (copies 1KB struct)
type LargeStruct struct {
    Data [128]int64
}

func process(s LargeStruct) {
    // ...
}
```

**Solution: Pass by reference**

```go
// Good: Passes pointer (8 bytes)
func process(s *LargeStruct) {
    // ...
}
```

### 6. Avoid Interface Conversions

**Problem: Unnecessary interface conversions**

```go
// Bad: Interface conversion in loop
func sumIntegers(values []interface{}) int {
    sum := 0
    for _, v := range values {
        sum += v.(int)
    }
    return sum
}
```

**Solution: Use concrete types**

```go
// Good: Direct access
func sumIntegers(values []int) int {
    sum := 0
    for _, v := range values {
        sum += v
    }
    return sum
}
```

### 7. Concurrent Processing

**Problem: Sequential processing**

```go
// Bad: Sequential
func processItems(items []Item) []Result {
    results := make([]Result, len(items))
    for i, item := range items {
        results[i] = process(item)
    }
    return results
}
```

**Solution: Parallel processing**

```go
// Good: Parallel
func processItems(items []Item) []Result {
    results := make([]Result, len(items))
    var wg sync.WaitGroup

    for i := range items {
        wg.Add(1)
        go func(idx int) {
            defer wg.Done()
            results[idx] = process(items[idx])
        }(i)
    }

    wg.Wait()
    return results
}
```

**With worker pool:**

```go
func processItems(items []Item) []Result {
    results := make([]Result, len(items))
    numWorkers := runtime.NumCPU()

    var wg sync.WaitGroup
    itemsChan := make(chan int, len(items))

    // Start workers
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for idx := range itemsChan {
                results[idx] = process(items[idx])
            }
        }()
    }

    // Send work
    for i := range items {
        itemsChan <- i
    }
    close(itemsChan)

    wg.Wait()
    return results
}
```

## Performance Testing

### Load Testing

**Create load test:**

```go
// test/load/load_test.go
//go:build load

package load_test

import (
    "context"
    "testing"
    "time"
)

func TestLoad(t *testing.T) {
    const (
        numRequests    = 10000
        numConcurrent  = 100
    )

    service := setupService()

    start := time.Now()

    sem := make(chan struct{}, numConcurrent)
    done := make(chan struct{})

    for i := 0; i < numRequests; i++ {
        sem <- struct{}{}
        go func() {
            defer func() {
                <-sem
                done <- struct{}{}
            }()

            _, err := service.Process(context.Background(), "test")
            if err != nil {
                t.Errorf("request failed: %v", err)
            }
        }()
    }

    for i := 0; i < numRequests; i++ {
        <-done
    }

    duration := time.Since(start)
    rps := float64(numRequests) / duration.Seconds()

    t.Logf("Duration: %v", duration)
    t.Logf("Requests/sec: %.2f", rps)
    t.Logf("Avg latency: %v", duration/time.Duration(numRequests))
}
```

**Run load test:**

```bash
go test -v -tags=load ./test/load -timeout 30m
```

### Stress Testing

```go
func TestStress(t *testing.T) {
    service := setupService()

    // Gradually increase load
    for concurrency := 10; concurrency <= 1000; concurrency *= 2 {
        t.Run(fmt.Sprintf("Concurrency-%d", concurrency), func(t *testing.T) {
            start := time.Now()

            var wg sync.WaitGroup
            for i := 0; i < concurrency; i++ {
                wg.Add(1)
                go func() {
                    defer wg.Done()
                    for j := 0; j < 100; j++ {
                        service.Process(context.Background(), "test")
                    }
                }()
            }

            wg.Wait()
            duration := time.Since(start)

            t.Logf("Concurrency %d: %v", concurrency, duration)
        })
    }
}
```

## Monitoring Performance

### Add Metrics

```go
import (
    "time"
    "log/slog"
)

type MetricsService struct {
    inner Service
}

func (m *MetricsService) Process(ctx context.Context, input string) (string, error) {
    start := time.Now()

    result, err := m.inner.Process(ctx, input)

    duration := time.Since(start)

    slog.InfoContext(ctx, "service.process",
        "duration_ms", duration.Milliseconds(),
        "input_size", len(input),
        "output_size", len(result),
        "error", err != nil,
    )

    return result, err
}
```

### Track Allocations

```go
import "runtime"

func trackAllocations(name string, fn func()) {
    var m1, m2 runtime.MemStats

    runtime.ReadMemStats(&m1)
    fn()
    runtime.ReadMemStats(&m2)

    slog.Info("allocation stats",
        "operation", name,
        "alloc_bytes", m2.TotalAlloc-m1.TotalAlloc,
        "num_allocs", m2.Mallocs-m1.Mallocs,
    )
}

// Usage
trackAllocations("process_data", func() {
    processData(data)
})
```

## Common Performance Issues

### 1. Goroutine Leaks

**Detect:**

```bash
# Add goroutine tracking
go test -v ./... -run TestMyService

# Check goroutine count
runtime.NumGoroutine()
```

**Fix:**

```go
// Bad: Goroutine never exits
func startWorker() {
    go func() {
        for {
            work := getWork()
            process(work)
        }
    }()
}

// Good: Use context for cancellation
func startWorker(ctx context.Context) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                work := getWork()
                process(work)
            }
        }
    }()
}
```

### 2. Memory Leaks

**Detect:**

```bash
# Memory profile over time
go tool pprof -alloc_space mem.prof
(pprof) top

# Look for growing allocations
```

**Fix:**

```go
// Bad: Holding references
var cache = make(map[string]*Data)

func store(key string, data *Data) {
    cache[key] = data  // Never cleaned up
}

// Good: Use LRU cache with eviction
type LRUCache struct {
    maxSize int
    items   map[string]*Data
}

func (c *LRUCache) Store(key string, data *Data) {
    if len(c.items) >= c.maxSize {
        c.evictOldest()
    }
    c.items[key] = data
}
```

### 3. Contention

**Detect:**

```bash
# Block profile
go tool pprof http://localhost:6060/debug/pprof/block

# Mutex profile
go tool pprof http://localhost:6060/debug/pprof/mutex
```

**Fix:**

```go
// Bad: Single mutex for everything
type Cache struct {
    mu    sync.Mutex
    items map[string]string
}

// Good: Sharded locks
type ShardedCache struct {
    shards [16]struct {
        mu    sync.RWMutex
        items map[string]string
    }
}

func (c *ShardedCache) getShard(key string) *struct {
    mu    sync.RWMutex
    items map[string]string
} {
    hash := fnv.New32a()
    hash.Write([]byte(key))
    return &c.shards[hash.Sum32()%16]
}
```

## Performance Checklist

Before releasing:

- [ ] Profile CPU usage (`-cpuprofile`)
- [ ] Profile memory allocations (`-memprofile`)
- [ ] Run benchmarks (`go test -bench`)
- [ ] Check for goroutine leaks
- [ ] Verify no memory leaks
- [ ] Test under load
- [ ] Measure startup time
- [ ] Check binary size
- [ ] Review hot paths
- [ ] Optimize allocations in loops

## Resources

- [Go Profiling](https://go.dev/blog/pprof)
- [Diagnostics](https://go.dev/doc/diagnostics)
- [Performance Tips](https://go.dev/wiki/Performance)
- [Memory Management](https://go.dev/blog/ismmkeynote)
- [Execution Tracer](https://go.dev/blog/execution-tracer)

## See Also

- **DEBUG_GUIDE.md** - Debugging techniques
- **TESTING_PATTERNS.md** - Benchmarking patterns
- **DOCKER_DEVELOPMENT.md** - Container optimization
