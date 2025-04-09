package version

import (
	"fmt"
	"io"
)

func Install(v string, w io.Writer) error {

	v = normaliseVersion(v)

	fmt.Fprintf(w, "Installing Go version %s\n", v)
	return nil
}

func Uninstall(v string, w io.Writer) error {

	v = normaliseVersion(v)

	fmt.Fprintf(w, "Uninstalling Go version %s\n", v)
	return nil
}
