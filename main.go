package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	args := os.Args[1:]

	if len(args) != 3 {
		fmt.Println("invalid arguments: usage: go run . <url> <max_concurrency> <max_pages>")
		os.Exit(1)
	}

	// Check if arguments provided for maxConcurrency & maxPages are integers
	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil || maxConcurrency <= 0 {
		fmt.Println("error: argument for maximum concurrency must be a positive integer")
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(args[2])
	if err != nil || maxPages <= 0 {
		fmt.Println("error: argument for maximum pages must be a postive integer")
		os.Exit(1)
	}

	rawBaseURL := args[0]

	parsedBaseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("invalid URL provided: %v\n", err)
		os.Exit(1)
	}

	pages := make(map[string]int)

	cfg := config{
		pages:              pages,
		baseURL:            parsedBaseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	fmt.Printf("starting crawl of: %s\n", rawBaseURL)

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	for normURL, count := range cfg.pages {
		fmt.Printf("%d - %s\n", count, normURL)
	}
}
