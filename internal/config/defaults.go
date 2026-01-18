package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// SetDefaults sets default values for all configuration options
func SetDefaults() {
	// General settings
	viper.SetDefault("verbose", false)
	viper.SetDefault("log_level", "info")

	// Output settings
	viper.SetDefault("output.format", "text")
	viper.SetDefault("output.color", true)
	viper.SetDefault("output.pretty", true)
	viper.SetDefault("output.quiet", false)
	viper.SetDefault("output.timestamp", false)
	viper.SetDefault("output.table_style", "ascii")

	// API settings
	viper.SetDefault("api.base_url", "https://api.example.com")
	viper.SetDefault("api.timeout", 30*time.Second)
	viper.SetDefault("api.retry_attempts", 3)
	viper.SetDefault("api.retry_delay", 1*time.Second)
	viper.SetDefault("api.follow_redirects", true)
	viper.SetDefault("api.verify_ssl", true)
	viper.SetDefault("api.user_agent", "termplate/1.0")
	viper.SetDefault("api.rate_limit_per_sec", 10)

	// Server settings
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 30*time.Second)
	viper.SetDefault("server.write_timeout", 30*time.Second)
	viper.SetDefault("server.idle_timeout", 60*time.Second)
	viper.SetDefault("server.shutdown_timeout", 10*time.Second)
	viper.SetDefault("server.tls_enabled", false)

	// File processing settings
	viper.SetDefault("files.input_dir", "./input")
	viper.SetDefault("files.output_dir", "./output")
	viper.SetDefault("files.temp_dir", getTempDir())
	viper.SetDefault("files.patterns", []string{"*"})
	viper.SetDefault("files.exclude_patterns", []string{})
	viper.SetDefault("files.max_file_size", 100*1024*1024) // 100MB
	viper.SetDefault("files.buffer_size", 4096)            // 4KB
	viper.SetDefault("files.create_dirs", true)
	viper.SetDefault("files.overwrite_existing", false)
	viper.SetDefault("files.preserve_perms", true)
	viper.SetDefault("files.backup_original", false)

	// Database settings
	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.database", "mydb")
	viper.SetDefault("database.username", "user")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", 5*time.Minute)
	viper.SetDefault("database.conn_max_idle_time", 10*time.Minute)
	viper.SetDefault("database.timeout", 10*time.Second)
	viper.SetDefault("database.migrations_path", "./migrations")
}

// getTempDir returns the system temp directory
func getTempDir() string {
	if tmpDir := os.Getenv("TMPDIR"); tmpDir != "" {
		return tmpDir
	}
	if tmpDir := os.TempDir(); tmpDir != "" {
		return tmpDir
	}
	return filepath.Join(".", "tmp")
}
