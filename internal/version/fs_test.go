package version

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOSFileSystem_Stat(t *testing.T) {
	fs := OSFileSystem{}

	// Test with a file that should exist
	_, err := fs.Stat("fs.go")
	if err != nil {
		t.Errorf("Expected fs.go to exist, got error: %v", err)
	}

	// Test with a file that shouldn't exist
	_, err = fs.Stat("nonexistent-file-1234567890.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestOSFileSystem_MkdirAndRemove(t *testing.T) {
	fs := OSFileSystem{}

	// Create a unique test directory
	testDir := filepath.Join(os.TempDir(), "gum-test-dir-"+filepath.Base(t.Name()))

	// Make sure it doesn't exist initially
	fs.RemoveAll(testDir)

	// Test MkdirAll
	err := fs.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Verify directory exists
	_, err = fs.Stat(testDir)
	if err != nil {
		t.Errorf("Test directory should exist after MkdirAll, got error: %v", err)
	}

	// Test RemoveAll
	err = fs.RemoveAll(testDir)
	if err != nil {
		t.Errorf("Failed to remove test directory: %v", err)
	}

	// Verify directory is gone
	_, err = fs.Stat(testDir)
	if !os.IsNotExist(err) {
		t.Errorf("Test directory should be gone after RemoveAll")
	}
}

func TestOSFileSystem_CreateTemp(t *testing.T) {
	fs := OSFileSystem{}

	// Test CreateTemp
	tempFile, err := fs.CreateTemp("", "gum-test-file-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Ensure we clean up
	tempName := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempName)

	// Verify file exists
	_, err = fs.Stat(tempName)
	if err != nil {
		t.Errorf("Temp file should exist after CreateTemp, got error: %v", err)
	}
}

func TestOSFileSystem_UserHomeDir(t *testing.T) {
	fs := OSFileSystem{}

	// Test UserHomeDir
	home, err := fs.UserHomeDir()
	if err != nil {
		t.Errorf("Failed to get user home directory: %v", err)
	}

	// Home directory should not be empty
	if home == "" {
		t.Error("Home directory should not be empty")
	}

	// Home directory should be a directory that exists
	fileInfo, err := fs.Stat(home)
	if err != nil {
		t.Errorf("Home directory does not exist: %v", err)
	} else if !fileInfo.IsDir() {
		t.Error("Home directory is not a directory")
	}
}
