package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(inputURL string) (string, error) {
	// Parse URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Combine host and path to get full path
	fullPath := parsedURL.Host + parsedURL.Path
	// Convert to lowercase
	fullPath = strings.ToLower(fullPath)
	// Remove trailing slash
	fullPath = strings.TrimSuffix(fullPath, "/")

	return fullPath, nil
}
