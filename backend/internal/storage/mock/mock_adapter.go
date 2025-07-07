package mock

import (
	"context"
	"io"
	"os"

	"github.com/ankuragarwal/fernfs/backend/internal/storage"
	"github.com/stretchr/testify/mock"
)

// Adapter is a mock implementation of storage.Adapter
type Adapter struct {
	mock.Mock
}

func (m *Adapter) Mkdir(ctx context.Context, path string, mode os.FileMode) error {
	args := m.Called(ctx, path, mode)
	return args.Error(0)
}

func (m *Adapter) Rmdir(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

func (m *Adapter) Readdir(ctx context.Context, path string) ([]storage.DirEntry, error) {
	args := m.Called(ctx, path)
	if entries := args.Get(0); entries != nil {
		return entries.([]storage.DirEntry), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *Adapter) Stat(ctx context.Context, path string) (*storage.FileInfo, error) {
	args := m.Called(ctx, path)
	if info := args.Get(0); info != nil {
		return info.(*storage.FileInfo), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *Adapter) ReadFile(ctx context.Context, path string, opts storage.ReadOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, path, opts)
	if reader := args.Get(0); reader != nil {
		return reader.(io.ReadCloser), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *Adapter) WriteFile(ctx context.Context, path string, data io.Reader, opts storage.WriteOptions) error {
	args := m.Called(ctx, path, data, opts)
	return args.Error(0)
}

func (m *Adapter) Unlink(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

func (m *Adapter) Rename(ctx context.Context, oldPath, newPath string) error {
	args := m.Called(ctx, oldPath, newPath)
	return args.Error(0)
}

func (m *Adapter) Chmod(ctx context.Context, path string, mode os.FileMode) error {
	args := m.Called(ctx, path, mode)
	return args.Error(0)
}

func (m *Adapter) Close() error {
	args := m.Called()
	return args.Error(0)
} 