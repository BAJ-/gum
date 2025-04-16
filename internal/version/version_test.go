package version

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

// MockFileSystem implements FileSystem for testing
type MockFileSystem struct {
	ExistingFiles map[string]bool
	DirError      error
	RemoveError   error
}

func (m *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	if m.ExistingFiles[name] {
		return nil, nil // Return a non-nil FileInfo for existing files
	}
	return nil, os.ErrNotExist
}

func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	if m.DirError != nil {
		return m.DirError
	}
	m.ExistingFiles[path] = true
	return nil
}

func (m *MockFileSystem) RemoveAll(path string) error {
	if m.RemoveError != nil {
		return m.RemoveError
	}
	delete(m.ExistingFiles, path)
	return nil
}

func (m *MockFileSystem) CreateTemp(dir, pattern string) (*os.File, error) {
	// Create a temp file in memory or return a mock file
	return os.CreateTemp("", pattern) // For simplicity, we'll use real temp files
}

func (m *MockFileSystem) UserHomeDir() (string, error) {
	return "/mock/home", nil
}

func TestVersionManager_Uninstall(t *testing.T) {
	tests := []struct {
		name         string
		version      string
		existingDirs map[string]bool
		removeError  error
		wantErr      bool
		wantOutput   string
	}{
		{
			name:         "uninstall success",
			version:      "go1.16.5",
			existingDirs: map[string]bool{"/mock/home/.gum/versions/go1.16.5": true},
			wantErr:      false,
			wantOutput:   "Successfully uninstalled",
		},
		{
			name:         "not installed",
			version:      "go1.16.5",
			existingDirs: map[string]bool{},
			wantErr:      false,
			wantOutput:   "not installed",
		},
		{
			name:         "remove error",
			version:      "go1.16.5",
			existingDirs: map[string]bool{"/mock/home/.gum/versions/go1.16.5": true},
			removeError:  errors.New("mock remove error"),
			wantErr:      true,
			wantOutput:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock filesystem
			mockFS := &MockFileSystem{
				ExistingFiles: tt.existingDirs,
				RemoveError:   tt.removeError,
			}

			// Create manager with mocks
			manager := &VersionManager{
				fs:         mockFS,
				installDir: "/mock/home/.gum/versions",
			}

			// Capture output
			var buf bytes.Buffer
			err := manager.Uninstall(tt.version, &buf)

			// Check error expectations
			if (err != nil) != tt.wantErr {
				t.Errorf("VersionManager.Uninstall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check output expectations
			if tt.wantOutput != "" && !strings.Contains(buf.String(), tt.wantOutput) {
				t.Errorf("Expected output to contain '%s', got '%s'", tt.wantOutput, buf.String())
			}
		})
	}
}

func TestExpandPath(t *testing.T) {
	mockFS := &MockFileSystem{}

	// Test with HOME variable
	path := "${HOME}/test/path"
	expanded := expandPath(path, mockFS)
	expected := "/mock/home/test/path"

	if expanded != expected {
		t.Errorf("expandPath(%s) = %s, want %s", path, expanded, expected)
	}

	// Test without HOME variable
	path = "/absolute/path"
	expanded = expandPath(path, mockFS)

	if expanded != path {
		t.Errorf("expandPath(%s) = %s, want %s", path, expanded, path)
	}
}
