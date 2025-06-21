package version

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestNormaliseVersion(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expected    string
		shouldPanic bool
	}{
		{
			name:        "already has go prefix",
			input:       "go1.24",
			expected:    "go1.24",
			shouldPanic: false,
		},
		{
			name:        "no go prefix",
			input:       "1.24",
			expected:    "go1.24",
			shouldPanic: false,
		},
		{
			name:        "with patch version",
			input:       "1.24.2",
			expected:    "go1.24.2",
			shouldPanic: false,
		},
		{
			name:        "with go prefix and patch version",
			input:       "go1.24.2",
			expected:    "go1.24.2",
			shouldPanic: false,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			shouldPanic: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("normaliseVersion(%q) did not panic as expected", tc.input)
					}
				}()
				normaliseVersion(tc.input)
			} else {
				result := normaliseVersion(tc.input)
				if result != tc.expected {
					t.Errorf("Expected normaliseVersion(%q) = %q,  got %q", tc.input, result, tc.expected)
				}
			}
		})
	}
}

func TestResolveVersion(t *testing.T) {
	// Mock HTTP client for testing
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			jsonContent := `[
				{
					"version": "go1.23.1",
					"stable": true,
					"files": []
				},
				{
					"version": "go1.23.2",
					"stable": true,
					"files": []
				},
				{
					"version": "go1.23.3",
					"stable": true,
					"files": []
				},
				{
					"version": "go1.24.1",
					"stable": true,
					"files": []
				},
				{
					"version": "go1.24.2",
					"stable": true,
					"files": []
				}
			]`

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(jsonContent)),
			}, nil
		},
	}

	testCases := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "complete version with go prefix",
			input:    "go1.23.4",
			expected: "go1.23.4",
			wantErr:  false,
		},
		{
			name:     "complete version without go prefix",
			input:    "1.23.4",
			expected: "1.23.4",
			wantErr:  false,
		},
		{
			name:     "major.minor version resolves to latest patch",
			input:    "1.23",
			expected: "go1.23.3",
			wantErr:  false,
		},
		{
			name:     "major.minor version with go prefix",
			input:    "go1.23",
			expected: "go1.23.3",
			wantErr:  false,
		},
		{
			name:     "invalid format returns as-is",
			input:    "invalid",
			expected: "invalid",
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := resolveVersion(tc.input, mockClient)

			if (err != nil) != tc.wantErr {
				t.Errorf("resolveVersion(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
				return
			}

			if result != tc.expected {
				t.Errorf("resolveVersion(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestResolveVersionHTTPError(t *testing.T) {
	// Mock HTTP client that returns an error
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	_, err := resolveVersion("1.23", mockClient)
	if err == nil {
		t.Error("Expected error when HTTP request fails, got nil")
	}

	if !strings.Contains(err.Error(), "failed to find latest patch version") {
		t.Errorf("Expected error to contain 'failed to find latest patch version', got: %v", err)
	}
}

func TestIsCompleteVersion(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"complete version", "1.23.4", true},
		{"major.minor", "1.23", false},
		{"single number", "1", false},
		{"invalid format", "1.23.4.5", false},
		{"with letters", "1.23a.4", false},
		{"empty string", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isCompleteVersion(tc.input)
			if result != tc.expected {
				t.Errorf("isCompleteVersion(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsMajorMinorVersion(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"major.minor", "1.23", true},
		{"complete version", "1.23.4", false},
		{"single number", "1", false},
		{"invalid format", "1.23.4.5", false},
		{"with letters", "1.23a", false},
		{"empty string", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isMajorMinorVersion(tc.input)
			if result != tc.expected {
				t.Errorf("isMajorMinorVersion(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestFetchAvailableVersions(t *testing.T) {
	// Mock HTTP client for testing
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.String(), "?mode=json") {
				t.Errorf("Expected request to include ?mode=json, got: %s", req.URL.String())
			}

			jsonContent := `[
				{
					"version": "go1.23.1",
					"stable": true,
					"files": []
				},
				{
					"version": "go1.23.2",
					"stable": true,
					"files": []
				},
				{
					"version": "go1.24.1",
					"stable": true,
					"files": []
				}
			]`

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(jsonContent)),
			}, nil
		},
	}

	versions, err := fetchAvailableVersions(mockClient)
	if err != nil {
		t.Fatalf("fetchAvailableVersions() error = %v", err)
	}

	expected := []string{"go1.23.1", "go1.23.2", "go1.24.1"}

	if len(versions) != len(expected) {
		t.Errorf("Expected %d versions, got %d", len(expected), len(versions))
	}

	for _, expectedVersion := range expected {
		found := false
		for _, version := range versions {
			if version == expectedVersion {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected version %s not found in result", expectedVersion)
		}
	}
}

func TestSemanticVersionSorting(t *testing.T) {
	// Mock HTTP client for testing
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			jsonContent := `[
				{
					"version": "go1.23.9",
					"stable": true,
					"files": []
				},
				{
					"version": "go1.23.10",
					"stable": true,
					"files": []
				},
				{
					"version": "go1.23.2",
					"stable": true,
					"files": []
				},
				{
					"version": "go1.24.1",
					"stable": true,
					"files": []
				}
			]`

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(jsonContent)),
			}, nil
		},
	}

	result, err := resolveVersion("1.23", mockClient)
	if err != nil {
		t.Fatalf("resolveVersion() error = %v", err)
	}

	expected := "go1.23.10"
	if result != expected {
		t.Errorf("resolveVersion('1.23') = %q, want %q", result, expected)
	}
}

func TestCompareVersions(t *testing.T) {
	testCases := []struct {
		name     string
		a        string
		b        string
		expected bool // a < b
	}{
		{"1.23.9 < 1.23.10", "go1.23.9", "go1.23.10", true},
		{"1.23.10 > 1.23.9", "go1.23.10", "go1.23.9", false},
		{"1.23.2 < 1.23.9", "go1.23.2", "go1.23.9", true},
		{"1.23.9 < 1.24.1", "go1.23.9", "go1.24.1", true},
		{"1.24.1 > 1.23.10", "go1.24.1", "go1.23.10", false},
		{"same versions", "go1.23.9", "go1.23.9", false},
		{"without go prefix", "1.23.9", "1.23.10", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := compareVersions(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("compareVersions(%q, %q) = %v, want %v", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}
