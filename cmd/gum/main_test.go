package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunCLI(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedErr    string
		expectedCode   int
	}{
		{
			name:         "no arguments",
			args:         []string{"gum"},
			expectedErr:  "Go Utility Manager (gum)",
			expectedCode: 1,
		},
		{
			name:           "install",
			args:           []string{"gum", "install", "1.24"},
			expectedOutput: "Installing Go version 1.24",
			expectedCode:   0,
		},
		{
			name:         "install without version",
			args:         []string{"gum", "install"},
			expectedErr:  "Error: no version provided",
			expectedCode: 1,
		},
		{
			name:           "uninstall",
			args:           []string{"gum", "uninstall", "1.24"},
			expectedOutput: "Uninstalling Go version 1.24",
			expectedCode:   0,
		},
		{
			name:         "uninstall without version",
			args:         []string{"gum", "uninstall"},
			expectedErr:  "Error: no version provided",
			expectedCode: 1,
		},
		{
			name:         "unknown command",
			args:         []string{"gum", "llatsni"},
			expectedErr:  "Unknown command: llatsni",
			expectedCode: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := runCLI(tc.args, &stdout, &stderr)

			if code != tc.expectedCode {
				t.Errorf("Expected exit code %d, got %d", tc.expectedCode, code)
			}

			if tc.expectedOutput != "" && !strings.Contains(stdout.String(), tc.expectedOutput) {
				t.Errorf("Expected stdout to contain '%s', got '%s'", tc.expectedOutput, stdout.String())
			}

			if tc.expectedErr != "" && !strings.Contains(stderr.String(), tc.expectedErr) {
				t.Errorf("Expected stderr to contain '%s', got '%s'", tc.expectedErr, stderr.String())
			}
		})
	}
}

func TestPrintUsage(t *testing.T) {
	var buf bytes.Buffer
	printUsage(&buf)

	expected := []string{
		"Go Utility Manager (gum)",
		"Usage:",
		"gum install <version>",
		"gum uninstall <version>",
	}

	for _, exp := range expected {
		if !strings.Contains(buf.String(), exp) {
			t.Errorf("Expected usage output to contain '%s', got '%s'", exp, buf.String())
		}
	}
}
