package config

import (
	"time"

	"github.com/spf13/viper"
)

func SetDefaults() {
	viper.SetDefault("verbose", false)
	viper.SetDefault("log_level", "info")
	viper.SetDefault("output", "text")
	viper.SetDefault("project.default_template", "default")
	viper.SetDefault("project.output_dir", ".")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 30*time.Second)
	viper.SetDefault("server.write_timeout", 30*time.Second)
	viper.SetDefault("api.base_url", "https://api.example.com")
	viper.SetDefault("api.timeout", 10*time.Second)
}
