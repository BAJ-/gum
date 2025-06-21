package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
)

func normaliseVersion(v string) string {
	if v == "" {
		panic("normaliseVersion: received empty string")
	}
	if !strings.HasPrefix(v, "go") {
		return "go" + v
	}
	return v
}

type progressWriter struct {
	w          io.Writer
	total      int64
	progress   int64
	startTime  time.Time
	lastUpdate time.Time
}

func newProgressWriter(w io.Writer, total int64) *progressWriter {
	return &progressWriter{
		w:          w,
		total:      total,
		startTime:  time.Now(),
		lastUpdate: time.Now(),
	}
}

// Write is required on the progress writer by io.TeeReader
func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.progress += int64(n)

	// Update every 100 ms
	if time.Since(pw.lastUpdate) < 100*time.Millisecond {
		return n, nil
	}

	pw.lastUpdate = time.Now()
	elapsed := time.Since(pw.startTime)

	if pw.total > 0 {
		pctDone := float64(pw.progress) / float64(pw.total) * 100
		fmt.Fprintf(pw.w, "\rDownloading... %.1f%% complete | %s elapsed", pctDone, formatDuration(elapsed))
	} else {
		fmt.Fprintf(pw.w, "\rDownloading... %s received in %s", formatBytes(pw.progress), formatDuration(elapsed))
	}

	return n, nil
}

func formatDuration(d time.Duration) string {
	if d.Hours() > 1 {
		return fmt.Sprintf("%.0fh %.0fm", d.Hours(), d.Minutes()-float64(int(d.Hours()))*60)
	} else if d.Minutes() > 1 {
		return fmt.Sprintf("%.0fm %.0fs", d.Minutes(), d.Seconds()-float64(int(d.Minutes()))*60)
	}
	return fmt.Sprintf("%.0fs", d.Seconds())
}

// formatBytes: Easier to read bytes format
func formatBytes(bytes int64) string {
	oneKB := int64(1024)
	// Less than 1 kilobyte
	if bytes < oneKB {
		return fmt.Sprintf("%d B", bytes)
	}

	// Up to 1 MB
	if bytes < oneKB*oneKB {
		kiloBytes := float64(bytes) / float64(oneKB)
		return fmt.Sprintf("%.1f KB", kiloBytes)
	}

	// Above 1MB
	megaBytes := float64(bytes) / (float64(oneKB) * float64(oneKB))
	return fmt.Sprintf("%.1f MB", megaBytes)
}

// resolveVersion resolves a version string to the full version
// If the version is already complete (e.g., "1.23.4"), it returns as-is
// If the version is major.minor (e.g., "1.23"), it finds the latest patch version
func resolveVersion(v string, client HTTPClient) (string, error) {
	cleanVersion := strings.TrimPrefix(v, "go")

	if isCompleteVersion(cleanVersion) {
		return v, nil
	}

	if isMajorMinorVersion(cleanVersion) {
		latestVersion, err := findLatestPatchVersion(cleanVersion, client)
		if err != nil {
			return "", fmt.Errorf("failed to find latest patch version for %s: %w", cleanVersion, err)
		}
		return latestVersion, nil
	}

	// If it doesn't match expected patterns, let the download fail naturally
	return v, nil
}

// isCompleteVersion checks if a version string is in major.minor.patch format
func isCompleteVersion(v string) bool {
	// Regex for major.minor.patch format (e.g., "1.23.4")
	completeVersionRegex := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	return completeVersionRegex.MatchString(v)
}

// isMajorMinorVersion checks if a version string is in major.minor format
func isMajorMinorVersion(v string) bool {
	// Regex for major.minor format (e.g., "1.23")
	majorMinorRegex := regexp.MustCompile(`^\d+\.\d+$`)
	return majorMinorRegex.MatchString(v)
}

// findLatestPatchVersion finds the latest patch version for a given major.minor version
func findLatestPatchVersion(majorMinor string, client HTTPClient) (string, error) {
	versions, err := fetchAvailableVersions(client)
	if err != nil {
		return "", fmt.Errorf("failed to fetch available versions: %w", err)
	}

	// Filter versions that match the major.minor pattern
	var matchingVersions []string
	prefix := "go" + majorMinor + "."

	for _, version := range versions {
		if strings.HasPrefix(version, prefix) {
			matchingVersions = append(matchingVersions, version)
		}
	}

	if len(matchingVersions) == 0 {
		return "", fmt.Errorf("no versions found for %s", majorMinor)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(matchingVersions)))

	return matchingVersions[0], nil
}

// GoVersion represents a Go version from the API
type GoVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []struct {
		Filename string `json:"filename"`
		OS       string `json:"os"`
		Arch     string `json:"arch"`
		Version  string `json:"version"`
		Kind     string `json:"kind"`
	} `json:"files"`
}

// fetchAvailableVersions fetches the list of available Go versions
func fetchAvailableVersions(client HTTPClient) ([]string, error) {
	req, err := http.NewRequest("GET", BaseURL+"/?mode=json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "gum/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch versions with status %s", resp.Status)
	}

	var versions []GoVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	var versionStrings []string
	for _, version := range versions {
		versionStrings = append(versionStrings, version.Version)
	}

	return versionStrings, nil
}
