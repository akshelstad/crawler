package main

import (
	"fmt"
	"sort"
)

type pagesReport struct {
	URL   string
	count int
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Println("=============================")
	fmt.Printf("REPORT for %s\n", baseURL)
	fmt.Println("=============================")

	pagesSorted := sortPages(pages)

	for _, page := range pagesSorted {
		fmt.Printf("Found %d internal links to %v\n", page.count, page.URL)
	}

}

func sortPages(pages map[string]int) []pagesReport {
	sortedPages := make([]pagesReport, 0, len(pages))
	for k, v := range pages {
		sortedPages = append(sortedPages, pagesReport{URL: k, count: v})
	}

	sort.Slice(sortedPages, func(i, j int) bool {
		return sortedPages[i].count > sortedPages[j].count
	})

	return sortedPages
}
