package main

import (
	"fmt"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) (map[string]int, error) {
	// Parse base URL
	parsedBaseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %v", err)
	}

	// Parse current URL
	parsedCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse current URL: %v", err)
	}

	// Check if current URL is from the same domain as base URL
	if parsedBaseURL.Host != parsedCurrentURL.Host {
		return nil, nil
	}

	// Normalize current URL
	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return nil, err
	}

	// Increment pages map
	pages[normalizedCurrentURL]++

	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		return nil, err
	}

	fmt.Println(pages)
	// fmt.Println(htmlBody)

	urls, err := getURLsFromHTML(htmlBody, rawCurrentURL)
	if err != nil {
		return nil, err
	}

	for _, u := range urls {
		if _, ok := pages[u]; !ok {
			pages, err = crawlPage(rawBaseURL, u, pages)
			if err != nil {
				return nil, err
			}
		}
	}

	return pages, nil
}
