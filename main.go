package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/pakkerman/web-crawler-go/graph"
)

func main() {
	start := time.Now()

	fmt.Println(len(os.Args))
	if len(os.Args) < 2 || 5 <= len(os.Args) {
		fmt.Println("# usage: ./crawler URL maxConcurrency maxPages")
		os.Exit(1)
	}

	parsedURL, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	maxConcurrency := 5
	if 2 < len(os.Args) && os.Args[2] != "" {
		val, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("# usage: ./crawler URL maxConcurrency maxPages")
			return
		}
		maxConcurrency = val
	}

	maxPages := 10
	if 3 < len(os.Args) && os.Args[3] != "" {
		val, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("# usage: ./crawler URL maxConcurrency maxPages")
			return
		}
		maxPages = val
	}

	cfg := config{
		pages:          make(map[string]int),
		externalLinks:  make(map[string]int),
		baseURL:        parsedURL,
		mu:             &sync.Mutex{},
		wg:             &sync.WaitGroup{},
		buffer:         make(chan struct{}, 10),
		counter:        0,
		maxConcurrency: maxConcurrency,
		maxPages:       maxPages,
	}

	fmt.Println("Concurrent requests: ", maxConcurrency)
	fmt.Println("Max Pages to crawl: ", maxPages)

	cfg.crawlPage(cfg.baseURL.String())
	cfg.wg.Wait()

	fmt.Println("=============================")
	fmt.Println("REPORT for", cfg.baseURL)
	fmt.Println("=============================")
	// cfg.printReport()

	fmt.Println("=============================")
	fmt.Printf("Pages crawled: %v\n", cfg.counter)
	fmt.Printf("Elapsed time: %s\n", time.Since(start))
	fmt.Println("=============================")

	fmt.Println("Links output to links.csv ")
	cfg.outputCSV()

	graph.RenderGraph()
}
