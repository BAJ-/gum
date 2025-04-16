package version

import "io"

// Manager defines the interface for version management operations
type Manager interface {
	Install(version string, w io.Writer) error
	Uninstall(version string, w io.Writer) error
}
