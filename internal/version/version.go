package version

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	InstallDir = "${HOME}/.gum/versions"
)

func Install(v string, w io.Writer) error {

	v = normaliseVersion(v)

	installDir := expandPath(InstallDir)
	versionDir := filepath.Join(installDir, v)

	if _, err := os.Stat(versionDir); err == nil {
		fmt.Fprintf(w, "Go %s is already installed at %s\n", v, versionDir)
		return nil
	}

	fmt.Fprintf(w, "Installing Go version %s\n", v)
	return nil
}

func Uninstall(v string, w io.Writer) error {

	v = normaliseVersion(v)

	fmt.Fprintf(w, "Uninstalling Go version %s\n", v)
	return nil
}
