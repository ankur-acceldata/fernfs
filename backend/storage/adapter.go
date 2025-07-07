package storage

import (
	"context"
	"io"
	"os"
	"time"
)

// FileInfo represents file metadata
type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}

// DirEntry represents a directory entry
type DirEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
}

// ReadOptions specifies options for reading files
type ReadOptions struct {
	Offset int64
	Length int64
}

// WriteOptions specifies options for writing files
type WriteOptions struct {
	Mode os.FileMode
}

// Adapter defines the interface for storage operations
type Adapter interface {
	// Mkdir creates a new directory
	Mkdir(ctx context.Context, path string, mode os.FileMode) error

	// Rmdir removes a directory
	Rmdir(ctx context.Context, path string) error

	// Readdir reads a directory
	Readdir(ctx context.Context, path string) ([]FileInfo, error)

	// Stat returns file information
	Stat(ctx context.Context, path string) (*FileInfo, error)

	// ReadFile reads a file
	ReadFile(ctx context.Context, path string, opts ReadOptions) (io.ReadCloser, error)

	// WriteFile writes a file
	WriteFile(ctx context.Context, path string, reader io.Reader, opts WriteOptions) error

	// Unlink removes a file
	Unlink(ctx context.Context, path string) error

	// Rename renames a file or directory
	Rename(ctx context.Context, oldPath, newPath string) error

	// Chmod changes file permissions
	Chmod(ctx context.Context, path string, mode os.FileMode) error
}

// Config represents the configuration for a storage adapter
type Config struct {
	Type string `mapstructure:"type"`
	// Common configuration options
	BasePath string `mapstructure:"base_path"`
	// Add more fields as needed for specific adapters
}

// Factory creates a new storage adapter instance
type Factory interface {
	Create(cfg Config) (Adapter, error)
} 