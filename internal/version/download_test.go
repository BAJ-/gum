package version

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		httpError  error
		httpStatus int
		httpBody   string
		wantErr    bool
	}{
		{
			name:       "successful download",
			url:        "https://example.com/test.tar.gz",
			httpStatus: http.StatusOK,
			httpBody:   "test file content",
			wantErr:    false,
		},
		{
			name:      "http error",
			url:       "https://example.com/test.tar.gz",
			httpError: errors.New("mock http error"),
			wantErr:   true,
		},
		{
			name:       "non-200 status",
			url:        "https://example.com/test.tar.gz",
			httpStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file for testing
			tmpFile, err := os.CreateTemp("", "download-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer tmpFile.Close()
			defer os.Remove(tmpFile.Name())

			// Set up mock HTTP client
			mockHTTP := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					// Verify request URL
					if req.URL.String() != tt.url {
						t.Errorf("Request URL = %v, want %v", req.URL.String(), tt.url)
					}

					// Verify User-Agent header
					if ua := req.Header.Get("User-Agent"); ua != "gum/1.0" {
						t.Errorf("User-Agent = %v, want gum/1.0", ua)
					}

					if tt.httpError != nil {
						return nil, tt.httpError
					}

					// Create a mock response
					body := io.NopCloser(bytes.NewBufferString(tt.httpBody))
					return &http.Response{
						StatusCode:    tt.httpStatus,
						Body:          body,
						ContentLength: int64(len(tt.httpBody)),
					}, nil
				},
			}

			// Capture output
			var buf bytes.Buffer
			err = downloadFile(tt.url, tmpFile, &buf, mockHTTP)

			// Check error expectations
			if (err != nil) != tt.wantErr {
				t.Errorf("downloadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// For successful downloads, check if the file contains the expected content
			if !tt.wantErr {
				// Reset file position
				if _, err := tmpFile.Seek(0, 0); err != nil {
					t.Fatalf("Failed to seek to beginning of file: %v", err)
				}

				content, err := io.ReadAll(tmpFile)
				if err != nil {
					t.Fatalf("Failed to read file content: %v", err)
				}

				if string(content) != tt.httpBody {
					t.Errorf("File content = %v, want %v", string(content), tt.httpBody)
				}

				// The progress info check isn't reliable in tests with small content
				// so we'll just verify some output was produced
				if buf.Len() == 0 {
					t.Errorf("Expected some output, but got none")
				}
			}
		})
	}
}
