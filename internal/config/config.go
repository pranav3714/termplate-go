package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Verbose  bool         `mapstructure:"verbose"`
	LogLevel string       `mapstructure:"log_level"`
	Output   OutputConfig `mapstructure:"output"`
	API      APIConfig    `mapstructure:"api"`
	Server   ServerConfig `mapstructure:"server"`
	Files    FilesConfig  `mapstructure:"files"`
	Database DBConfig     `mapstructure:"database"`
}

// OutputConfig controls output formatting
type OutputConfig struct {
	Format      string `mapstructure:"format"`      // text, json, yaml, table, csv
	ColorOutput bool   `mapstructure:"color"`       // Enable colored output
	Pretty      bool   `mapstructure:"pretty"`      // Pretty print JSON/YAML
	Quiet       bool   `mapstructure:"quiet"`       // Minimal output
	Timestamp   bool   `mapstructure:"timestamp"`   // Include timestamps
	TableStyle  string `mapstructure:"table_style"` // ascii, unicode, markdown
}

// APIConfig holds API client configuration
type APIConfig struct {
	BaseURL         string            `mapstructure:"base_url"`
	Key             string            `mapstructure:"key"`
	Secret          string            `mapstructure:"secret"`
	Token           string            `mapstructure:"token"`
	Timeout         time.Duration     `mapstructure:"timeout"`
	RetryAttempts   int               `mapstructure:"retry_attempts"`
	RetryDelay      time.Duration     `mapstructure:"retry_delay"`
	FollowRedirects bool              `mapstructure:"follow_redirects"`
	VerifySSL       bool              `mapstructure:"verify_ssl"`
	UserAgent       string            `mapstructure:"user_agent"`
	Headers         map[string]string `mapstructure:"headers"`
	RateLimitPerSec int               `mapstructure:"rate_limit_per_sec"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	TLSEnabled      bool          `mapstructure:"tls_enabled"`
	TLSCertFile     string        `mapstructure:"tls_cert_file"`
	TLSKeyFile      string        `mapstructure:"tls_key_file"`
}

// FilesConfig holds file processing configuration
type FilesConfig struct {
	InputDir          string   `mapstructure:"input_dir"`
	OutputDir         string   `mapstructure:"output_dir"`
	TempDir           string   `mapstructure:"temp_dir"`
	Patterns          []string `mapstructure:"patterns"`           // File patterns to match
	ExcludePatterns   []string `mapstructure:"exclude_patterns"`   // Patterns to exclude
	MaxFileSize       int64    `mapstructure:"max_file_size"`      // Max file size in bytes
	BufferSize        int      `mapstructure:"buffer_size"`        // Buffer size for reading
	CreateDirs        bool     `mapstructure:"create_dirs"`        // Auto-create directories
	OverwriteExisting bool     `mapstructure:"overwrite_existing"` // Overwrite existing files
	PreservePerms     bool     `mapstructure:"preserve_perms"`     // Preserve file permissions
	BackupOriginal    bool     `mapstructure:"backup_original"`    // Backup before processing
}

// DBConfig holds database configuration
type DBConfig struct {
	Driver          string        `mapstructure:"driver"` // postgres, mysql, sqlite
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Database        string        `mapstructure:"database"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	SSLMode         string        `mapstructure:"ssl_mode"` // disable, require, verify-ca, verify-full
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	Timeout         time.Duration `mapstructure:"timeout"`
	MigrationsPath  string        `mapstructure:"migrations_path"`
}

// Load reads configuration from viper
func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}
	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate output format
	validFormats := map[string]bool{
		"text": true, "json": true, "yaml": true, "table": true, "csv": true,
	}
	if !validFormats[c.Output.Format] {
		return fmt.Errorf("invalid output format: %s (valid: text, json, yaml, table, csv)", c.Output.Format)
	}

	// Validate server port
	if c.Server.Port < 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	// Validate database port if driver is specified
	if c.Database.Driver != "" && (c.Database.Port < 0 || c.Database.Port > 65535) {
		return fmt.Errorf("invalid database port: %d", c.Database.Port)
	}

	// Validate file size limit
	if c.Files.MaxFileSize < 0 {
		return fmt.Errorf("invalid max file size: %d", c.Files.MaxFileSize)
	}

	// Validate API retry attempts
	if c.API.RetryAttempts < 0 {
		return fmt.Errorf("invalid retry attempts: %d", c.API.RetryAttempts)
	}

	return nil
}

// GetAPIAuthHeader returns the appropriate authorization header
func (c *APIConfig) GetAPIAuthHeader() (string, string) {
	if c.Token != "" {
		return "Authorization", fmt.Sprintf("Bearer %s", c.Token)
	}
	if c.Key != "" {
		return "X-API-Key", c.Key
	}
	return "", ""
}

// GetDSN returns the database connection string
func (c *DBConfig) GetDSN() string {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			c.Username, c.Password, c.Host, c.Port, c.Database)
	case "sqlite":
		return c.Database // SQLite uses file path as DSN
	default:
		return ""
	}
}
