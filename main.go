package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

type config struct {
	pages          map[string]int
	baseURL        *url.URL
	mu             *sync.Mutex
	wg             *sync.WaitGroup
	buffer         chan struct{}
	counter        int
	maxPages       int
	maxConcurrency int
}

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
	fmt.Println("REPORT for ", cfg.baseURL)
	fmt.Println("=============================")
	fmt.Printf("Pages crawled: %v\n", cfg.counter)
	fmt.Printf("Elapsed time: %s\n", time.Since(start))
	cfg.printReport()
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.wg.Add(1)
	defer func() {
		cfg.wg.Done()
	}()

	fmt.Println("crawling page:", rawCurrentURL)

	body, err := getHTML(rawCurrentURL)
	if err != nil {
		// fmt.Println(err)
		return
	}

	urls, _ := getURLsFromHTML(body, rawCurrentURL)

	var unvisited []string
	for _, link := range urls {
		// fmt.Println("\tProcessing: ", link)

		parsedLink, _ := url.Parse(link)
		if parsedLink.Host != cfg.baseURL.Host {
			continue
		}

		normalizedURL, err := normalizeURL(link)
		if err != nil {
			fmt.Println(err)
			continue
		}

		isFirst := cfg.addPageVisit(normalizedURL)
		if !isFirst {
			continue
		}

		unvisited = append(unvisited, link)
	}

	// fmt.Println("\tunvisited urls on this page: ", len(unvisited))

	for _, link := range unvisited {
		cfg.incrementCounter()
		fmt.Println(cfg.counter)
		go func(link string) {
			cfg.buffer <- struct{}{}
			cfg.crawlPage(link)
			<-cfg.buffer
		}(link)

		if cfg.pagesMaxed() {
			break
		}
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	_, ok := cfg.pages[normalizedURL]
	if ok {
		cfg.pages[normalizedURL] += 1
		return false
	}

	cfg.pages[normalizedURL] = 1
	return true
}

func (cfg *config) incrementCounter() {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.counter++
}

func (cfg config) pagesMaxed() bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	return cfg.maxPages <= cfg.counter
}

func (cfg config) printReport() {
	for key, value := range cfg.pages {
		fmt.Println(key, value)
	}
}
