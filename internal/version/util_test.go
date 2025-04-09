package version

import "testing"

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
