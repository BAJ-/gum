package version

import (
	"os"
)

// FileSystem abstracts file system operations for better testability
type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
	RemoveAll(path string) error
	CreateTemp(dir, pattern string) (*os.File, error)
	UserHomeDir() (string, error)
}

// OSFileSystem implements FileSystem using the os package
type OSFileSystem struct{}

func (fs OSFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (fs OSFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs OSFileSystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (fs OSFileSystem) CreateTemp(dir, pattern string) (*os.File, error) {
	return os.CreateTemp(dir, pattern)
}

func (fs OSFileSystem) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}
