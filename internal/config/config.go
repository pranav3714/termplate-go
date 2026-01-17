package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Verbose  bool          `mapstructure:"verbose"`
	LogLevel string        `mapstructure:"log_level"`
	Output   string        `mapstructure:"output"`
	Project  ProjectConfig `mapstructure:"project"`
	Server   ServerConfig  `mapstructure:"server"`
	API      APIConfig     `mapstructure:"api"`
}

type ProjectConfig struct {
	DefaultTemplate string `mapstructure:"default_template"`
	OutputDir       string `mapstructure:"output_dir"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type APIConfig struct {
	BaseURL string        `mapstructure:"base_url"`
	Key     string        `mapstructure:"key"`
	Timeout time.Duration `mapstructure:"timeout"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Server.Port < 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid port: %d", c.Server.Port)
	}
	return nil
}
