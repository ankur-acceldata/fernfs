package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
}

type LoggingConfig struct {
	Level string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "15s")
	v.SetDefault("server.shutdown_timeout", "5s")
	v.SetDefault("metrics.enabled", true)
	v.SetDefault("metrics.path", "/metrics")
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")

	// Configure environment variables
	v.SetEnvPrefix("FERNFS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Bind environment variables
	bindEnv(v, "server.host")
	bindEnv(v, "server.port")
	bindEnv(v, "server.read_timeout")
	bindEnv(v, "server.write_timeout")
	bindEnv(v, "server.shutdown_timeout")
	bindEnv(v, "metrics.enabled")
	bindEnv(v, "metrics.path")
	bindEnv(v, "logging.level")
	bindEnv(v, "logging.format")

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func bindEnv(v *viper.Viper, key string) {
	if err := v.BindEnv(key, fmt.Sprintf("FERNFS_%s", strings.ToUpper(strings.ReplaceAll(key, ".", "_")))); err != nil {
		panic(fmt.Sprintf("failed to bind env var for key %s: %v", key, err))
	}
} 