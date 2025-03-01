package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

const maxConcurrency = 5

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		return
	}
	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		return
	}

	rawBaseURL := args[0]

	parsedBaseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("invalid URL provided: %v\n", err)
		return
	}

	pages := make(map[string]int)

	cfg := config{
		pages:              pages,
		baseURL:            parsedBaseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
	}

	fmt.Printf("starting crawl of: %s\n", rawBaseURL)

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	for normURL, count := range cfg.pages {
		fmt.Printf("%d - %s\n", count, normURL)
	}
}
