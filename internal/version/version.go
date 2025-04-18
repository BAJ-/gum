package version

import (
	"bufio"
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
	httpClient HTTPClient
	installDir string
}

// NewManager creates a new Manager with default implementations
func NewManager() Manager {
	return &VersionManager{
		fs:         OSFileSystem{},
		httpClient: NewDefaultHTTPClient(),
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
	if err := downloadAndExtract(downloadURL, versionDir, w, m.httpClient); err != nil {
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

// Use creates a symlink to make the specified Go version active
func (m *VersionManager) Use(v string, w io.Writer) error {
	if v == "" {
		goModVersion, err := detectVersionInGoMod(m.fs)
		if err != nil {
			return fmt.Errorf("Failed to detect version in go.mod: %w", err)
		}
		v = goModVersion
		fmt.Fprintf(w, "Detected Go %s from go.mod\n", v)
	}

	v = normaliseVersion(v)
	versionDir := filepath.Join(m.installDir, v)

	// Check if version is installed
	if _, err := m.fs.Stat(versionDir); os.IsNotExist(err) {
		return fmt.Errorf("Go %s is not installed. Use 'gum install %s' first", v, v)
	}

	// Get user home directory
	home, err := m.fs.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create .gum/bin directory if it doesn't already exist
	// This is where active go versions will be linked from
	binDir := filepath.Join(home, ".gum", "bin")
	if err := m.fs.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Find target go binary for the symlink
	srcPath := filepath.Join(versionDir, "bin", "go")
	if _, err := m.fs.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("Go binary not found in %s", versionDir)
	}

	// Path to the symlink
	linkPath := filepath.Join(binDir, "go")

	// Check if symlink already exists and remove it
	if _, err := m.fs.Stat(linkPath); err == nil {
		// Try reading the current link to see where it points
		currentTarget, err := m.fs.ReadLink(linkPath)
		if err == nil && filepath.Clean(currentTarget) == filepath.Clean(srcPath) {
			// Symlink already points to requested version
			fmt.Fprintf(w, "Go %s is already the active version\n", v)
			return nil
		}

		// Symlink point to wrong version so we remove it
		if err := m.fs.Remove(linkPath); err != nil {
			return fmt.Errorf("failed to update Go version: %w", err)
		}
	}

	// Create new symlink pointing to requested version
	if err := m.fs.Symlink(srcPath, linkPath); err != nil {
		return fmt.Errorf("failed to set Go %s as active: %w", v, err)
	}

	fmt.Fprintf(w, "Successfully set Go %s as the active version\n", v)

	return nil
}

func (m *VersionManager) List(w io.Writer) error {
	entries, err := os.ReadDir(m.installDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(w, "No Go versions installed yet")
			return nil
		}

		return fmt.Errorf("failed to read versions directory: %w", err)
	}

	if len(entries) == 0 {
		fmt.Fprintln(w, "No Go versions installed yet")
		return nil
	}

	// Find active version, if any
	activeVersion := ""
	home, _ := m.fs.UserHomeDir()
	binDir := filepath.Join(home, ".gum", "bin")
	// Full path to the 'go' command symlink
	linkPath := filepath.Join(binDir, "go")

	resolvedPath, err := filepath.EvalSymlinks(linkPath)
	if err == nil {
		// The version binary is stored two folders deep in the version folder,
		// so we need to go two folders up to get the version
		versionDir := filepath.Dir(filepath.Dir(resolvedPath))
		// Then we can get the version from the folder name
		activeVersion = filepath.Base(versionDir)
	}

	fmt.Fprintln(w, "Installed Go versions:")

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		version := entry.Name()
		if version == activeVersion {
			fmt.Fprintf(w, "* %s (active)\n", version)
		} else {
			fmt.Fprintf(w, "  %s\n", version)
		}
	}
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

// detectVersionInGoMod read go.mod in current directory
// and tries to extract the Go version
func detectVersionInGoMod(fs FileSystem) (string, error) {

	if _, err := fs.Stat("go.mod"); os.IsNotExist(err) {
		return "", nil
	}

	modFile, err := os.Open("go.mod")
	if err != nil {
		return "", fmt.Errorf("could not open go.mod file: %w", err)
	}
	defer modFile.Close()

	// Read the file line by line to find the go version
	scanner := bufio.NewScanner(modFile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Look for lines that start with 'go'
		if strings.HasPrefix(line, "go ") {
			version := strings.TrimPrefix(line, "go ")
			// Assume the first line we find that starts with 'go'
			// Is the one that defines the Go version
			if version == "" {
				return "", fmt.Errorf("invalid Go version format in go.mod")
			}
			return version, nil
		}
	}

	return "", fmt.Errorf("no Go version found in go.mod")
}
