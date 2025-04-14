package version

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	defaultInstallDir = "${HOME}/.gum/versions"
)

func Install(v string, w io.Writer) error {

	v = normaliseVersion(v)

	installDir := expandPath(defaultInstallDir)
	versionDir := filepath.Join(installDir, v)

	if _, err := os.Stat(versionDir); err == nil {
		fmt.Fprintf(w, "Go %s is already installed at %s\n", v, versionDir)
		return nil
	}

	// Get download url based on version and architecture
	downloadURL, err := getDownloadURL(v)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("failed to create installation directory: %w", err)
	}

	fmt.Fprintf(w, "Downloading %s...\n", downloadURL)
	if err := downloadAndExtract(downloadURL, versionDir, w); err != nil {
		os.RemoveAll(versionDir)
		return err
	}

	fmt.Fprintf(w, "Successfully installed Go %s at %s\n", v, versionDir)
	return nil
}

func Uninstall(v string, w io.Writer) error {

	v = normaliseVersion(v)

	fmt.Fprintf(w, "Uninstalling Go version %s\n", v)
	return nil
}
