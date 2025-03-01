package main

import (
	"fmt"
	"net/url"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	// Parse current URL
	parsedCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("error - crawlPage: failed to parse current URL: %v\n", err)
		return
	}

	// Check if current URL is from the same domain as base URL
	if cfg.baseURL.Host != parsedCurrentURL.Host {
		return
	}

	// Normalize current URL
	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error - normalizeURL: %v\n", err)
		return
	}

	// Check if page has already been crawled and increment
	if !cfg.addPageVisit(normalizedCurrentURL) {
		return
	}

	fmt.Printf("crawling %s\n", rawCurrentURL)

	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error - getHTML: %v\n", err)
		return
	}

	urls, err := getURLsFromHTML(htmlBody, cfg.baseURL)
	if err != nil {
		fmt.Printf("error - getURLsFromHTML: %v\n", err)
		return
	}

	for _, url := range urls {
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}
}

func (cfg *config) addPageVisit(normURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, ok := cfg.pages[normURL]; ok {
		cfg.pages[normURL]++
		return false
	}

	cfg.pages[normURL] = 1
	return true
}
