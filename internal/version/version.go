package version

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultInstallDir = "${HOME}/.gum/versions"
)

// VersionManager handles Go version installation and uninstallation
type VersionManager struct {
	fs         FileSystem
	installDir string
}

// NewManager creates a new Manager with default implementations
func NewManager() Manager {
	return &VersionManager{
		fs:         OSFileSystem{},
		installDir: expandPath(defaultInstallDir, OSFileSystem{}),
	}
}

// Install installs a specific Go version
func (m *VersionManager) Install(v string, w io.Writer) error {
	v = normaliseVersion(v)
	versionDir := filepath.Join(m.installDir, v)

	// Check if already installed
	if _, err := m.fs.Stat(versionDir); err == nil {
		fmt.Fprintf(w, "Go %s is already installed at %s\n", v, versionDir)
		return nil
	}

	// Get download url based on version and architecture
	downloadURL, err := getDownloadURL(v)
	if err != nil {
		return err
	}

	if err := m.fs.MkdirAll(m.installDir, 0755); err != nil {
		return fmt.Errorf("failed to create installation directory: %w", err)
	}

	fmt.Fprintf(w, "Downloading %s...\n", downloadURL)
	if err := downloadAndExtract(downloadURL, versionDir, w); err != nil {
		m.fs.RemoveAll(versionDir)
		return err
	}

	fmt.Fprintf(w, "Successfully installed Go %s at %s\n", v, versionDir)
	return nil
}

// Uninstall removes a specific Go version
func (m *VersionManager) Uninstall(v string, w io.Writer) error {
	v = normaliseVersion(v)
	versionDir := filepath.Join(m.installDir, v)

	fmt.Fprintf(w, "Uninstalling Go version %s\n", v)

	// Check if the version is installed
	if _, err := m.fs.Stat(versionDir); os.IsNotExist(err) {
		fmt.Fprintf(w, "Go %s is not installed at %s\n", v, versionDir)
		return nil
	}

	// Remove the version directory
	if err := m.fs.RemoveAll(versionDir); err != nil {
		return fmt.Errorf("failed to uninstall Go %s: %w", v, err)
	}

	fmt.Fprintf(w, "Successfully uninstalled Go %s from %s\n", v, versionDir)
	return nil
}

// Utility function to expand paths using the filesystem
func expandPath(path string, fs FileSystem) string {
	if strings.HasPrefix(path, "${HOME}") {
		home, err := fs.UserHomeDir()
		if err == nil {
			return strings.Replace(path, "${HOME}", home, 1)
		}
	}
	return path
}
