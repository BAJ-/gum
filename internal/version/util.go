package version

import (
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
