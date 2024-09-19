package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"
)

type config struct {
	pages   map[string]int
	baseURL *url.URL
	mu      *sync.Mutex
	wg      *sync.WaitGroup
	buffer  chan struct{}
	counter int
}

func main() {
	start := time.Now()

	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)

	} else if 2 < len(os.Args) {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	parsedURL, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("starting crawl of:", parsedURL.String())

	cfg := config{
		pages:   make(map[string]int),
		baseURL: parsedURL,
		mu:      &sync.Mutex{},
		wg:      &sync.WaitGroup{},
		buffer:  make(chan struct{}, 10),
		counter: 0,
	}

	cfg.crawlPage(cfg.baseURL.String())
	cfg.wg.Wait()

	fmt.Printf("Pages crawled: %v\n", cfg.counter)
	fmt.Printf("Elapsed time: %s\n", time.Since(start))
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.wg.Add(1)
	defer func() {
		cfg.wg.Done()
		cfg.incrementCounter()
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

	if len(unvisited) == 0 {
		return
	}

	fmt.Println("\tunvisited urls on this page: ", len(unvisited))

	for _, link := range unvisited {
		go func(link string) {
			cfg.buffer <- struct{}{}
			cfg.crawlPage(link)
			<-cfg.buffer
		}(link)
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
