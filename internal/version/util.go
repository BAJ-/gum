package version

import (
	"os"
	"strings"
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

func expandPath(path string) string {
	if strings.HasPrefix(path, "${HOME}") {
		home, err := os.UserHomeDir()
		if err == nil {
			return strings.Replace(path, "${HOME}", home, 1)
		}
	}
	return path
}
