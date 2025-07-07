package local

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ankuragarwal/fernfs/backend/internal/storage"
	"github.com/ankuragarwal/fernfs/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdapter(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := testutils.CreateTempDir(t)

	// Create adapter instance
	adapter, err := New(Config{BasePath: tempDir})
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("directory operations", func(t *testing.T) {
		// Test Mkdir
		err := adapter.Mkdir(ctx, "testdir", 0755)
		require.NoError(t, err)

		// Verify directory exists
		info, err := os.Stat(filepath.Join(tempDir, "testdir"))
		require.NoError(t, err)
		assert.True(t, info.IsDir())
		assert.Equal(t, os.FileMode(0755), info.Mode().Perm())

		// Test Readdir on empty directory
		entries, err := adapter.Readdir(ctx, "testdir")
		require.NoError(t, err)
		assert.Empty(t, entries)

		// Test Rmdir
		err = adapter.Rmdir(ctx, "testdir")
		require.NoError(t, err)

		// Verify directory is gone
		_, err = os.Stat(filepath.Join(tempDir, "testdir"))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("file operations", func(t *testing.T) {
		// Test WriteFile
		content := []byte("test content")
		err := adapter.WriteFile(ctx, "test.txt", bytes.NewReader(content), storage.WriteOptions{Mode: 0644})
		require.NoError(t, err)

		// Test Stat
		info, err := adapter.Stat(ctx, "test.txt")
		require.NoError(t, err)
		assert.Equal(t, "test.txt", info.Name)
		assert.Equal(t, int64(len(content)), info.Size)
		assert.Equal(t, os.FileMode(0644), info.Mode.Perm())
		assert.False(t, info.IsDir)

		// Test ReadFile
		reader, err := adapter.ReadFile(ctx, "test.txt", storage.ReadOptions{})
		require.NoError(t, err)
		data, err := io.ReadAll(reader)
		require.NoError(t, err)
		reader.Close()
		assert.Equal(t, content, data)

		// Test ReadFile with offset and length
		reader, err = adapter.ReadFile(ctx, "test.txt", storage.ReadOptions{Offset: 5, Length: 4})
		require.NoError(t, err)
		data, err = io.ReadAll(reader)
		require.NoError(t, err)
		reader.Close()
		assert.Equal(t, []byte("cont"), data)

		// Test Chmod
		err = adapter.Chmod(ctx, "test.txt", 0600)
		require.NoError(t, err)
		info, err = adapter.Stat(ctx, "test.txt")
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0600), info.Mode.Perm())

		// Test Rename
		err = adapter.Rename(ctx, "test.txt", "renamed.txt")
		require.NoError(t, err)
		_, err = adapter.Stat(ctx, "test.txt")
		assert.Error(t, err)
		info, err = adapter.Stat(ctx, "renamed.txt")
		require.NoError(t, err)
		assert.Equal(t, "renamed.txt", info.Name)

		// Test Unlink
		err = adapter.Unlink(ctx, "renamed.txt")
		require.NoError(t, err)
		_, err = adapter.Stat(ctx, "renamed.txt")
		assert.Error(t, err)
	})

	t.Run("path validation", func(t *testing.T) {
		// Test absolute path
		err := adapter.Mkdir(ctx, "/absolute/path", 0755)
		assert.Error(t, err)

		// Test path traversal
		err = adapter.Mkdir(ctx, "../outside", 0755)
		assert.Error(t, err)

		// Test empty path
		err = adapter.Mkdir(ctx, "", 0755)
		assert.Error(t, err)
	})

	t.Run("error cases", func(t *testing.T) {
		// Test Rmdir on non-existent directory
		err := adapter.Rmdir(ctx, "nonexistent")
		assert.Error(t, err)

		// Test Rmdir on file
		err = adapter.WriteFile(ctx, "test.txt", bytes.NewReader([]byte("test")), storage.WriteOptions{Mode: 0644})
		require.NoError(t, err)
		err = adapter.Rmdir(ctx, "test.txt")
		assert.Error(t, err)

		// Test Unlink on directory
		err = adapter.Mkdir(ctx, "testdir", 0755)
		require.NoError(t, err)
		err = adapter.Unlink(ctx, "testdir")
		assert.Error(t, err)

		// Clean up
		adapter.Unlink(ctx, "test.txt")
		adapter.Rmdir(ctx, "testdir")
	})
}

func TestNew(t *testing.T) {
	t.Run("empty base path", func(t *testing.T) {
		_, err := New(Config{BasePath: ""})
		assert.Error(t, err)
	})

	t.Run("invalid base path", func(t *testing.T) {
		_, err := New(Config{BasePath: string([]byte{0})})
		assert.Error(t, err)
	})

	t.Run("successful creation", func(t *testing.T) {
		tempDir := testutils.CreateTempDir(t)
		adapter, err := New(Config{BasePath: tempDir})
		require.NoError(t, err)
		assert.NotNil(t, adapter)
		assert.Equal(t, tempDir, adapter.basePath)
	})
} 