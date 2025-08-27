package fs

import (
	"context"
	"os"
)

// fileSystem implements FileSystem
type fileSystem struct{}

// NewFileSystem creates a new filesystem
func NewFileSystem() FileSystem {
	return &fileSystem{}
}

func (fs *fileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (fs *fileSystem) WriteFile(path string, data []byte, perm int) error {
	return os.WriteFile(path, data, os.FileMode(perm))
}

func (fs *fileSystem) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (fs *fileSystem) MkdirAll(path string, perm int) error {
	return os.MkdirAll(path, os.FileMode(perm))
}

// configLoader implements ConfigLoader
type configLoader struct {
	fs FileSystem
}

// NewConfigLoader creates a new config loader
func NewConfigLoader(fs FileSystem) ConfigLoader {
	return &configLoader{fs: fs}
}

func (cl *configLoader) Load(ctx context.Context, path string) ([]byte, error) {
	return cl.fs.ReadFile(path)
}
