package testutils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ankuragarwal/fernfs/backend/pkg/config"
	"github.com/stretchr/testify/require"
)

// GetTestConfig returns a test configuration
func GetTestConfig(t *testing.T) *config.Config {
	t.Helper()
	return &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8081, // Different from main port to avoid conflicts
		},
		Metrics: config.MetricsConfig{
			Enabled: true,
			Path:    "/metrics",
		},
		Logging: config.LoggingConfig{
			Level:  "debug",
			Format: "console",
		},
	}
}

// CreateTempDir creates a temporary directory and returns its path.
// The directory will be automatically cleaned up when the test finishes.
func CreateTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "fernfs-test-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	return dir
}

// GetTestDataPath returns the absolute path to a test data file
func GetTestDataPath(t *testing.T, relativePath string) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Join(filepath.Dir(filename), "../../test/data", relativePath)
}

// CreateTestFile creates a test file with the given content
func CreateTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)
	return path
} 