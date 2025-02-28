package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// getURLsFromHTML extracts URLs from an HTML body and applies a base URL to relative URLs
func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	// Check if base URL is empty
	if rawBaseURL == "" {
		return nil, fmt.Errorf("base URL is empty")
	}

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %v", err)
	}

	htmlReader := strings.NewReader(htmlBody)
	doc, err := html.Parse(htmlReader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var urls []string
	var traverseNodes func(*html.Node)
	traverseNodes = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, anchor := range node.Attr {
				if anchor.Key == "href" {
					href, err := url.Parse(anchor.Val)
					if err != nil {
						fmt.Printf("failed to parse href '%v': %v\n", anchor.Val, err)
						continue
					}

					resolvedURL := baseURL.ResolveReference(href)
					urls = append(urls, resolvedURL.String())
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverseNodes(child)
		}
	}
	traverseNodes(doc)

	return urls, nil
}

func getHTML(rawURL string) (string, error) {
	// Fetch webpage at rawURL
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 { // Check for 400+ status code
		return "", fmt.Errorf("got HTTP error: %v", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") { // Check if content type is HTML
		return "", fmt.Errorf("invalid content type: %v", contentType)
	}

	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Convert HTML bytes to string
	htmlBody := string(htmlBytes)

	return htmlBody, nil
}
