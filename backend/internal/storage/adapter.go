package storage

import (
	"context"
	"io"
	"os"
	"time"
)

// FileInfo represents metadata about a file
type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime time.Time `json:"mod_time"`
	IsDir   bool      `json:"is_dir"`
}

// DirEntry represents a directory entry
type DirEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
}

// WriteOptions represents options for write operations
type WriteOptions struct {
	Mode os.FileMode
	// Add more options as needed (e.g., metadata, encryption)
}

// ReadOptions represents options for read operations
type ReadOptions struct {
	Offset int64
	Length int64
	// Add more options as needed (e.g., decryption)
}

// Adapter defines the interface that all storage backends must implement
type Adapter interface {
	// Directory operations
	Mkdir(ctx context.Context, path string, mode os.FileMode) error
	Rmdir(ctx context.Context, path string) error
	Readdir(ctx context.Context, path string) ([]DirEntry, error)

	// File operations
	Stat(ctx context.Context, path string) (*FileInfo, error)
	ReadFile(ctx context.Context, path string, opts ReadOptions) (io.ReadCloser, error)
	WriteFile(ctx context.Context, path string, data io.Reader, opts WriteOptions) error
	Unlink(ctx context.Context, path string) error
	Rename(ctx context.Context, oldPath, newPath string) error
	Chmod(ctx context.Context, path string, mode os.FileMode) error

	// Lifecycle operations
	Close() error
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