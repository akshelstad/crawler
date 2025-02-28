package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputBody string
		inputURL  string
		expected  []string
		expectErr string
	}{
		{
			name:      "no URLs",
			inputBody: "<html></html>",
			inputURL:  "https://example.com",
			expected:  []string{},
		},
		{
			name: "single URL",
			inputBody: `
			<html>
				<a href="/path/one"></a>
			</html>
			`,
			inputURL: "https://example.com",
			expected: []string{"https://example.com/path/one"},
		},
		{
			name: "relative and absolute URLs",
			inputBody: `
			<html>
				<body>	
					<a href="/path/one">
						<span>example.com</span>
					</a>
					<a href="https://other.com/path/one">
						<span>example.com</span>
					</a>
				</body>		
			</html>
			`,
			inputURL: "https://example.com",
			expected: []string{"https://example.com/path/one", "https://other.com/path/one"},
		},
		{
			name: "URL with full URL",
			inputBody: `
			<html>
				<body>
					<a href="https://example.com/path/two">
						<span>example.com</span> 
					</a>
				</body>
			</html>
			`,
			inputURL: "https://example.com",
			expected: []string{"https://example.com/path/two"},
		},
		{
			name: "no input URL",
			inputBody: `
			<html>
				<body>
					<a href="/path/three">
						<span>example.com</span>
					</a>
				</body>
			</html>
			`,
			inputURL:  "",
			expected:  []string{},
			expectErr: "base URL is empty",
		},
		{
			name: "invalid href URL",
			inputBody: `
			<html>
				<body>
					<a href=":\\invalid">
						<span>example.com</span>
					</a>
				</body>
			</html>
			`,
			inputURL:  "https://example.com",
			expected:  []string{},
			expectErr: "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil && !strings.Contains(err.Error(), tc.expectErr) {
				t.Errorf("Test %v - %s FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err != nil && tc.expectErr == "" {
				t.Errorf("Test %v - %s FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err == nil && tc.expectErr != "" {
				t.Errorf("Test %v - %s FAIL: expected error: %v, actual: %v", i, tc.name, tc.expectErr, err)
				return
			}

			if len(actual) != len(tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected %v URLs, actual: %v", i, tc.name, len(tc.expected), len(actual))
				return
			}

			for j, url := range actual {
				if !reflect.DeepEqual(url, tc.expected[j]) {
					t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected[j], url)
				}
			}
		})
	}

}
