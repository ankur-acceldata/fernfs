package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ankuragarwal/fernfs/backend/internal/storage"
)

// Adapter implements storage.Adapter for local filesystem
type Adapter struct {
	basePath string
}

// Config represents local filesystem adapter configuration
type Config struct {
	BasePath string `mapstructure:"base_path"`
}

// New creates a new local filesystem adapter
func New(cfg Config) (*Adapter, error) {
	if cfg.BasePath == "" {
		return nil, fmt.Errorf("base_path is required")
	}

	// Create base directory if it doesn't exist
	if err := os.MkdirAll(cfg.BasePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	return &Adapter{
		basePath: cfg.BasePath,
	}, nil
}

// resolvePath resolves a relative path to an absolute path within the base directory
func (a *Adapter) resolvePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Clean the path to remove any . or .. components
	cleanPath := filepath.Clean(path)
	if cleanPath == ".." || filepath.IsAbs(cleanPath) {
		return "", fmt.Errorf("invalid path: %s", path)
	}

	// Join with base path and verify it's within the base directory
	fullPath := filepath.Join(a.basePath, cleanPath)
	if !filepath.HasPrefix(fullPath, a.basePath) {
		return "", fmt.Errorf("path escapes base directory: %s", path)
	}

	return fullPath, nil
}

// Mkdir creates a new directory
func (a *Adapter) Mkdir(ctx context.Context, path string, mode os.FileMode) error {
	fullPath, err := a.resolvePath(path)
	if err != nil {
		return err
	}
	return os.MkdirAll(fullPath, mode)
}

// Rmdir removes a directory
func (a *Adapter) Rmdir(ctx context.Context, path string) error {
	fullPath, err := a.resolvePath(path)
	if err != nil {
		return err
	}

	// Verify it's a directory
	info, err := os.Stat(fullPath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", path)
	}

	return os.Remove(fullPath)
}

// Readdir lists directory contents
func (a *Adapter) Readdir(ctx context.Context, path string) ([]storage.DirEntry, error) {
	fullPath, err := a.resolvePath(path)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	result := make([]storage.DirEntry, len(entries))
	for i, entry := range entries {
		result[i] = storage.DirEntry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
		}
	}

	return result, nil
}

// Stat returns file information
func (a *Adapter) Stat(ctx context.Context, path string) (*storage.FileInfo, error) {
	fullPath, err := a.resolvePath(path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	return &storage.FileInfo{
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
	}, nil
}

// ReadFile reads a file's contents
func (a *Adapter) ReadFile(ctx context.Context, path string, opts storage.ReadOptions) (io.ReadCloser, error) {
	fullPath, err := a.resolvePath(path)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(fullPath, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	if opts.Offset > 0 {
		if _, err := file.Seek(opts.Offset, io.SeekStart); err != nil {
			file.Close()
			return nil, err
		}
	}

	if opts.Length > 0 {
		return &limitedReadCloser{
			ReadCloser: file,
			limit:      opts.Length,
		}, nil
	}

	return file, nil
}

// WriteFile writes data to a file
func (a *Adapter) WriteFile(ctx context.Context, path string, data io.Reader, opts storage.WriteOptions) error {
	fullPath, err := a.resolvePath(path)
	if err != nil {
		return err
	}

	// Create parent directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, opts.Mode)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	return err
}

// Unlink removes a file
func (a *Adapter) Unlink(ctx context.Context, path string) error {
	fullPath, err := a.resolvePath(path)
	if err != nil {
		return err
	}

	// Verify it's not a directory
	info, err := os.Stat(fullPath)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("cannot unlink directory: %s", path)
	}

	return os.Remove(fullPath)
}

// Rename moves/renames a file or directory
func (a *Adapter) Rename(ctx context.Context, oldPath, newPath string) error {
	oldFullPath, err := a.resolvePath(oldPath)
	if err != nil {
		return err
	}

	newFullPath, err := a.resolvePath(newPath)
	if err != nil {
		return err
	}

	// Create parent directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(newFullPath), 0755); err != nil {
		return err
	}

	return os.Rename(oldFullPath, newFullPath)
}

// Chmod changes file permissions
func (a *Adapter) Chmod(ctx context.Context, path string, mode os.FileMode) error {
	fullPath, err := a.resolvePath(path)
	if err != nil {
		return err
	}

	return os.Chmod(fullPath, mode)
}

// Close implements storage.Adapter
func (a *Adapter) Close() error {
	return nil
}

// limitedReadCloser wraps an io.ReadCloser to limit the number of bytes that can be read
type limitedReadCloser struct {
	io.ReadCloser
	limit int64
	read  int64
}

func (l *limitedReadCloser) Read(p []byte) (n int, err error) {
	if l.limit <= l.read {
		return 0, io.EOF
	}
	if int64(len(p)) > l.limit-l.read {
		p = p[0 : l.limit-l.read]
	}
	n, err = l.ReadCloser.Read(p)
	l.read += int64(n)
	return
} 