# Configuration Guide

Your CLI now has comprehensive configuration support for API clients, file processing, database connections, and multiple output formats.

## Table of Contents

- [Configuration Files](#configuration-files)
- [Environment Variables](#environment-variables)
- [Configuration Structure](#configuration-structure)
- [Using Configuration in Code](#using-configuration-in-code)
- [Output Formatting](#output-formatting)
- [Examples](#examples)

## Configuration Files

Configuration can be loaded from:

1. **Command line flag**: `--config /path/to/config.yaml`
2. **Home directory**: `~/.ever-so-powerful-go.yaml`
3. **Current directory**: `./.ever-so-powerful-go.yaml`

### Create Your Config File

```bash
# Copy the example configuration
cp configs/config.example.yaml ~/.ever-so-powerful-go.yaml

# Edit with your settings
vim ~/.ever-so-powerful-go.yaml
```

## Environment Variables

All configuration can be overridden with environment variables using the prefix `EVER_SO_POWERFUL_GO_`:

```bash
# General settings
export EVER_SO_POWERFUL_GO_VERBOSE=true
export EVER_SO_POWERFUL_GO_LOG_LEVEL=debug

# Output settings
export EVER_SO_POWERFUL_GO_OUTPUT_FORMAT=json
export EVER_SO_POWERFUL_GO_OUTPUT_COLOR=false

# API settings
export EVER_SO_POWERFUL_GO_API_KEY=your-api-key
export EVER_SO_POWERFUL_GO_API_BASE_URL=https://api.example.com
export EVER_SO_POWERFUL_GO_API_TIMEOUT=60s

# Database settings
export EVER_SO_POWERFUL_GO_DB_USER=dbuser
export EVER_SO_POWERFUL_GO_DB_PASSWORD=dbpass
export EVER_SO_POWERFUL_GO_DB_HOST=localhost
export EVER_SO_POWERFUL_GO_DB_PORT=5432

# File processing
export EVER_SO_POWERFUL_GO_FILES_INPUT_DIR=/path/to/input
export EVER_SO_POWERFUL_GO_FILES_OUTPUT_DIR=/path/to/output
```

## Configuration Structure

### General Settings

```yaml
verbose: false
log_level: info  # debug, info, warn, error
```

### Output Configuration

```yaml
output:
  format: text          # text, json, yaml, table, csv
  color: true           # Enable colored output
  pretty: true          # Pretty print JSON/YAML
  quiet: false          # Minimal output
  timestamp: false      # Include timestamps
  table_style: ascii    # ascii, unicode, markdown
```

### API Configuration

```yaml
api:
  base_url: https://api.example.com
  key: ${EVER_SO_POWERFUL_GO_API_KEY}
  token: ${EVER_SO_POWERFUL_GO_API_TOKEN}
  timeout: 30s
  retry_attempts: 3
  retry_delay: 1s
  follow_redirects: true
  verify_ssl: true
  user_agent: "ever-so-powerful-go/1.0"
  headers:
    X-Custom-Header: "value"
  rate_limit_per_sec: 10
```

### Server Configuration

```yaml
server:
  host: localhost
  port: 8080
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 60s
  shutdown_timeout: 10s
  tls_enabled: false
  tls_cert_file: /path/to/cert.pem
  tls_key_file: /path/to/key.pem
```

### File Processing Configuration

```yaml
files:
  input_dir: ./input
  output_dir: ./output
  temp_dir: /tmp
  patterns:
    - "*.txt"
    - "*.json"
  exclude_patterns:
    - "*.tmp"
    - ".*"
  max_file_size: 104857600  # 100MB
  buffer_size: 4096
  create_dirs: true
  overwrite_existing: false
  preserve_perms: true
  backup_original: false
```

### Database Configuration

```yaml
database:
  driver: postgres  # postgres, mysql, sqlite
  host: localhost
  port: 5432
  database: mydb
  username: ${EVER_SO_POWERFUL_GO_DB_USER}
  password: ${EVER_SO_POWERFUL_GO_DB_PASSWORD}
  ssl_mode: disable
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m
  conn_max_idle_time: 10m
  timeout: 10s
  migrations_path: ./migrations
```

## Using Configuration in Code

### Loading Configuration

```go
package mycommand

import (
    "github.com/spf13/viper"
    "github.com/blacksilver/ever-so-powerful/internal/config"
)

func runCommand() error {
    // Load full configuration
    cfg, err := config.Load()
    if err != nil {
        return fmt.Errorf("loading config: %w", err)
    }

    // Validate configuration
    if err := cfg.Validate(); err != nil {
        return fmt.Errorf("invalid config: %w", err)
    }

    // Use configuration
    fmt.Printf("API Base URL: %s\n", cfg.API.BaseURL)
    fmt.Printf("Output Format: %s\n", cfg.Output.Format)
}
```

### Accessing Individual Settings

```go
// Using viper directly
apiKey := viper.GetString("api.key")
timeout := viper.GetDuration("api.timeout")
retries := viper.GetInt("api.retry_attempts")
verbose := viper.GetBool("verbose")

// File settings
inputDir := viper.GetString("files.input_dir")
patterns := viper.GetStringSlice("files.patterns")
maxSize := viper.GetInt64("files.max_file_size")

// Database settings
dbHost := viper.GetString("database.host")
dbPort := viper.GetInt("database.port")
```

### Using Helper Methods

```go
// Get API auth header
cfg, _ := config.Load()
headerName, headerValue := cfg.API.GetAPIAuthHeader()
// Returns: "Authorization", "Bearer <token>"
// or: "X-API-Key", "<api-key>"

// Get database connection string
dsn := cfg.Database.GetDSN()
// Returns: "host=localhost port=5432 user=... dbname=..."
```

## Output Formatting

Use the output formatter to easily support multiple output formats:

### Basic Usage

```go
package handler

import (
    "github.com/spf13/viper"
    "github.com/blacksilver/ever-so-powerful/internal/config"
    "github.com/blacksilver/ever-so-powerful/internal/output"
)

func Execute(ctx context.Context, in Input) error {
    // Load output config
    cfg, _ := config.Load()
    formatter := output.NewFormatter(cfg.Output)

    // Prepare data
    result := map[string]string{
        "name":   "example",
        "status": "success",
        "count":  "42",
    }

    // Print in configured format (text, json, yaml, table, csv)
    return formatter.Print(result)
}
```

### Different Data Types

#### Simple Map

```go
data := map[string]string{
    "key1": "value1",
    "key2": "value2",
}
formatter.Print(data)

// Output (table format):
// | Key  | Value  |
// |------|--------|
// | key1 | value1 |
// | key2 | value2 |
```

#### Slice of Maps (Table Data)

```go
data := []map[string]string{
    {"name": "Alice", "age": "30", "city": "NYC"},
    {"name": "Bob", "age": "25", "city": "LA"},
}
formatter.Print(data)

// Output (table format):
// | name  | age | city |
// |-------|-----|------|
// | Alice | 30  | NYC  |
// | Bob   | 25  | LA   |
```

#### 2D String Array

```go
data := [][]string{
    {"Name", "Age", "City"},
    {"Alice", "30", "NYC"},
    {"Bob", "25", "LA"},
}
formatter.Print(data)
```

### Table Styles

Set the table style in config or via environment variable:

```bash
# ASCII style (default)
export EVER_SO_POWERFUL_GO_OUTPUT_TABLE_STYLE=ascii

# Unicode style (pretty box characters)
export EVER_SO_POWERFUL_GO_OUTPUT_TABLE_STYLE=unicode

# Markdown style
export EVER_SO_POWERFUL_GO_OUTPUT_TABLE_STYLE=markdown
```

## Examples

### Example 1: API Client Configuration

```yaml
# ~/.ever-so-powerful-go.yaml
api:
  base_url: https://api.github.com
  token: ${GITHUB_TOKEN}
  timeout: 60s
  retry_attempts: 3
  headers:
    Accept: "application/vnd.github.v3+json"
```

```go
// In your code
cfg, _ := config.Load()

// Create HTTP client
client := &http.Client{
    Timeout: cfg.API.Timeout,
}

// Create request
req, _ := http.NewRequest("GET", cfg.API.BaseURL+"/user", nil)

// Add auth header
headerName, headerValue := cfg.API.GetAPIAuthHeader()
if headerName != "" {
    req.Header.Set(headerName, headerValue)
}

// Add custom headers
for k, v := range cfg.API.Headers {
    req.Header.Set(k, v)
}
```

### Example 2: File Processing

```yaml
files:
  input_dir: ./data/input
  output_dir: ./data/output
  patterns:
    - "*.csv"
    - "*.json"
  max_file_size: 52428800  # 50MB
  create_dirs: true
  backup_original: true
```

```go
cfg, _ := config.Load()

// Ensure directories exist
if cfg.Files.CreateDirs {
    os.MkdirAll(cfg.Files.InputDir, 0755)
    os.MkdirAll(cfg.Files.OutputDir, 0755)
}

// Process files
files, _ := filepath.Glob(filepath.Join(cfg.Files.InputDir, cfg.Files.Patterns[0]))
for _, file := range files {
    // Check file size
    info, _ := os.Stat(file)
    if info.Size() > cfg.Files.MaxFileSize {
        continue // Skip large files
    }

    // Backup if configured
    if cfg.Files.BackupOriginal {
        backupPath := file + ".bak"
        io.Copy(...)
    }

    // Process file...
}
```

### Example 3: Database Connection

```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  database: myapp
  username: ${DB_USER}
  password: ${DB_PASSWORD}
  max_open_conns: 25
  conn_max_lifetime: 5m
```

```go
cfg, _ := config.Load()

// Get connection string
dsn := cfg.Database.GetDSN()

// Open database
db, err := sql.Open(cfg.Database.Driver, dsn)
if err != nil {
    return err
}

// Configure connection pool
db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
db.SetConnMaxIdleTime(cfg.Database.ConnMaxIdleTime)
```

### Example 4: Output Formatting

```go
// Configure output format
cfg := config.OutputConfig{
    Format:     "json",  // or "yaml", "table", "csv"
    Pretty:     true,
    ColorOutput: true,
}

formatter := output.NewFormatter(cfg)

// Print structured data
data := []map[string]string{
    {"id": "1", "name": "Item 1", "status": "active"},
    {"id": "2", "name": "Item 2", "status": "pending"},
}

// Will output as JSON, YAML, table, or CSV based on config
formatter.Print(data)
```

## Testing Configuration

Create a test configuration file:

```bash
# Create test config
cat > test-config.yaml << EOF
verbose: true
log_level: debug
output:
  format: json
  pretty: true
api:
  base_url: http://localhost:8080
  timeout: 5s
EOF

# Run with test config
./build/bin/ever-so-powerful-go --config test-config.yaml example greet --name "Test"
```

## Configuration Priority

Configuration values are resolved in this order (highest to lowest priority):

1. **Command-line flags**: `--verbose`
2. **Environment variables**: `EVER_SO_POWERFUL_GO_VERBOSE=true`
3. **Config file**: `verbose: true` in `~/.ever-so-powerful-go.yaml`
4. **Defaults**: Set in `internal/config/defaults.go`

## Best Practices

1. **Use environment variables for secrets**: Never commit API keys or passwords
2. **Create environment-specific configs**: `config.dev.yaml`, `config.prod.yaml`
3. **Validate configuration**: Always call `cfg.Validate()` after loading
4. **Document your settings**: Add comments to config files
5. **Use helper methods**: Leverage `GetAPIAuthHeader()`, `GetDSN()`, etc.

## Next Steps

- See example usage in `cmd/example/greet.go`
- Check the full config example in `configs/config.example.yaml`
- Read about adding commands in `docs/GETTING_STARTED.md`
