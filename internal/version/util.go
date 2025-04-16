package version

import (
	"fmt"
	"io"
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
