package version

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	BaseURL = "https://golang.org/dl"
)

func getDownloadURL(v string) (string, error) {
	installOS := runtime.GOOS
	arch := runtime.GOARCH

	// Not all architectures returned by runtime match format
	// used for Go version downloads
	archMap := map[string]string{
		"amd64": "amd64",
		"386":   "386",
		"arm64": "arm64",
		"arm":   "armv6l",
	}

	archName, ok := archMap[arch]
	if !ok {
		panic(fmt.Sprintf("unsupported architecture %s.", arch))
	}

	switch installOS {
	case "darwin":
		if arch == "amd64" || arch == "arm64" {
			return fmt.Sprintf("%s/%s.darwin-%s.tar.gz", BaseURL, v, archName), nil
		}
		return "", fmt.Errorf("unsupported macOS architecture: %s", archName)
	case "linux":
		return fmt.Sprintf("%s/%s.%s-%s.tar.gz", BaseURL, v, installOS, archName), nil
	default:
		return "", fmt.Errorf("OS not supported by gum: %s", installOS)
	}
}

func downloadAndExtract(url, destDir string, w io.Writer, client HTTPClient) error {
	tmpFile, err := os.CreateTemp("", "gum-download-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	// Clean up temp file deferred
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	if err := downloadFile(url, tmpFile, w, client); err != nil {
		return err
	}

	// Move pointer to beginning of file
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to prepare for extraction: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(destDir), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	fmt.Fprintf(w, "Extracting to %s...\n", destDir)
	if strings.HasSuffix(url, ".tar.gz") {
		err = extractTarGz(tmpFile, destDir)
	} else {
		err = fmt.Errorf("unsupported archive format: %s", url)
	}

	if err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	return nil
}

func downloadFile(url string, file *os.File, w io.Writer, client HTTPClient) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "gum/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %s", resp.Status)
	}

	progress := newProgressWriter(w, resp.ContentLength)

	_, err = io.Copy(file, io.TeeReader(resp.Body, progress))
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// New line
	fmt.Fprintln(w)

	return nil
}
