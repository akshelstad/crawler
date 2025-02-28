package main

import (
	"strings"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		expected  string
		expectErr string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://example.com/path",
			expected: "example.com/path",
		},
		{
			name:     "remove trailing slash",
			inputURL: "https://example.com/path/",
			expected: "example.com/path",
		},
		{
			name:     "convert to lowercase",
			inputURL: "https://eXample.com/PATH",
			expected: "example.com/path",
		},
		{
			name:     "remove scheme, trailing slash, and convert to lowercase",
			inputURL: "https://eXample.com/PATH/",
			expected: "example.com/path",
		},
		{
			name:      "invalid URL",
			inputURL:  ":\\invalid",
			expected:  "",
			expectErr: "failed to parse URL",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil && !strings.Contains(err.Error(), tc.expectErr) {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err != nil && tc.expectErr == "" {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err == nil && tc.expectErr != "" {
				t.Errorf("Test %v - '%s' FAIL: expected error: %v, actual: %v", i, tc.name, tc.expectErr, err)
				return
			}

			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
