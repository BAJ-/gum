package main

import (
	"fmt"
	"io"
	"os"

	"github.com/baj-/gum/internal/version"
)

var versionManager version.Manager = version.NewManager()

func main() {
	os.Exit(runCLI(os.Args, os.Stdout, os.Stderr))
}

func runCLI(args []string, stdout, stderr io.Writer) int {
	if len(args) < 2 {
		printUsage(stderr)
		return 1
	}

	command := args[1]

	switch command {
	case "install":
		if len(args) < 3 {
			fmt.Fprintln(stderr, "Error: no version provided")
			printUsage(stderr)
			return 1
		}
		versionStr := args[2]
		err := versionManager.Install(versionStr, stdout)
		if err != nil {
			fmt.Fprintf(stderr, "Error installing Go %s: %v\n", versionStr, err)
			return 1
		}
		return 0
	case "uninstall":
		if len(args) < 3 {
			fmt.Fprintln(stderr, "Error: no version provided")
			printUsage(stderr)
			return 1
		}
		versionStr := args[2]
		err := versionManager.Uninstall(versionStr, stdout)
		if err != nil {
			fmt.Fprintf(stderr, "Error uninstalling Go %s: %v\n", versionStr, err)
			return 1
		}
		return 0
	case "use":
		versionStr := ""
		if len(args) >= 3 {
			versionStr = args[2]
		}

		err := versionManager.Use(versionStr, stdout)
		if err != nil {
			fmt.Fprintf(stderr, "Error setting Go %s as active: %v\n", versionStr, err)
			return 1
		}
		return 0
	case "list":
		err := versionManager.List(stdout)
		if err != nil {
			fmt.Fprintf(stderr, "Error listing Go versions: %v\n", err)
			return 1
		}
		return 0
	default:
		fmt.Fprintf(stderr, "Unknown command: %s\n", command)
		printUsage(stderr)
		return 1
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Go Utility Manager (gum)")
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  gum install <version>   - Install a specific Go version")
	fmt.Fprintln(w, "  gum uninstall <version> - Uninstall a specific Go version")
	fmt.Fprintln(w, "  gum use <version>       - Use a specific Go version")
	fmt.Fprintln(w, "  gum list                - List installed versions")
}
