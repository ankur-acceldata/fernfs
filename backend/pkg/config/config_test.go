package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		wantConfig  *Config
		wantErr     bool
		description string
	}{
		{
			name: "default config",
			envVars: map[string]string{},
			wantConfig: &Config{
				Server: ServerConfig{
					Host:            "0.0.0.0",
					Port:            8080,
					ReadTimeout:     15 * time.Second,
					WriteTimeout:    15 * time.Second,
					ShutdownTimeout: 5 * time.Second,
				},
				Metrics: MetricsConfig{
					Enabled: true,
					Path:    "/metrics",
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "json",
				},
			},
			wantErr:     false,
			description: "Should load default configuration when no environment variables are set",
		},
		{
			name: "custom config from env",
			envVars: map[string]string{
				"FERNFS_SERVER_HOST":             "127.0.0.1",
				"FERNFS_SERVER_PORT":             "9090",
				"FERNFS_SERVER_READ_TIMEOUT":     "30s",
				"FERNFS_SERVER_WRITE_TIMEOUT":    "30s",
				"FERNFS_SERVER_SHUTDOWN_TIMEOUT": "10s",
				"FERNFS_METRICS_ENABLED":         "false",
				"FERNFS_METRICS_PATH":            "/custom-metrics",
				"FERNFS_LOGGING_LEVEL":           "debug",
				"FERNFS_LOGGING_FORMAT":          "console",
			},
			wantConfig: &Config{
				Server: ServerConfig{
					Host:            "127.0.0.1",
					Port:            9090,
					ReadTimeout:     30 * time.Second,
					WriteTimeout:    30 * time.Second,
					ShutdownTimeout: 10 * time.Second,
				},
				Metrics: MetricsConfig{
					Enabled: false,
					Path:    "/custom-metrics",
				},
				Logging: LoggingConfig{
					Level:  "debug",
					Format: "console",
				},
			},
			wantErr:     false,
			description: "Should override default configuration with environment variables",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment variables
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			// Load configuration
			got, err := LoadConfig()

			// Check error
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Compare configuration
			assert.Equal(t, tt.wantConfig.Server.Host, got.Server.Host)
			assert.Equal(t, tt.wantConfig.Server.Port, got.Server.Port)
			assert.Equal(t, tt.wantConfig.Server.ReadTimeout, got.Server.ReadTimeout)
			assert.Equal(t, tt.wantConfig.Server.WriteTimeout, got.Server.WriteTimeout)
			assert.Equal(t, tt.wantConfig.Server.ShutdownTimeout, got.Server.ShutdownTimeout)
			assert.Equal(t, tt.wantConfig.Metrics.Enabled, got.Metrics.Enabled)
			assert.Equal(t, tt.wantConfig.Metrics.Path, got.Metrics.Path)
			assert.Equal(t, tt.wantConfig.Logging.Level, got.Logging.Level)
			assert.Equal(t, tt.wantConfig.Logging.Format, got.Logging.Format)
		})
	}
} 