package version

import "io"

type Manager interface {
	Install(version string, w io.Writer) error
	Uninstall(version string, w io.Writer) error
}

type DefaultManager struct{}

func NewManager() Manager {
	return &DefaultManager{}
}

func (m *DefaultManager) Install(version string, w io.Writer) error {
	return Install(version, w)
}

func (m *DefaultManager) Uninstall(version string, w io.Writer) error {
	return Uninstall(version, w)
}
