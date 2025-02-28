package main

import (
	"fmt"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	// Parse base URL
	parsedBaseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("error - crawlPage: failed to parse base URL: %v\n", err)
		return
	}

	// Parse current URL
	parsedCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("error - crawlPage: failed to parse current URL: %v\n", err)
		return
	}

	// Check if current URL is from the same domain as base URL
	if parsedBaseURL.Host != parsedCurrentURL.Host {
		return
	}

	// Normalize current URL
	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error - normalizeURL: %v\n", err)
		return
	}

	// Check if page has already been crawled and increment
	if _, ok := pages[normalizedCurrentURL]; ok {
		pages[normalizedCurrentURL]++
		return
	}

	// Mark as visited
	pages[normalizedCurrentURL] = 1

	fmt.Printf("crawling %s\n", rawCurrentURL)

	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error - getHTML: %v\n", err)
		return
	}

	urls, err := getURLsFromHTML(htmlBody, rawBaseURL)
	if err != nil {
		fmt.Printf("error - getURLsFromHTML: %v\n", err)
		return
	}

	for _, url := range urls {
		crawlPage(rawBaseURL, url, pages)
	}
}
