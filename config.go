package main

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/pakkerman/web-crawler-go/helpers/links"
)

type config struct {
	pages          map[string]int
	externalLinks  map[string]int
	baseURL        *url.URL
	mu             *sync.Mutex
	wg             *sync.WaitGroup
	buffer         chan struct{}
	counter        int
	maxPages       int
	maxConcurrency int
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.wg.Add(1)
	defer func() {
		cfg.wg.Done()
	}()

	fmt.Println("crawling page:", rawCurrentURL)

	body, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}

	urls, _ := getURLsFromHTML(body, rawCurrentURL)

	var unvisited []string
	for _, link := range urls {
		// fmt.Println("\tProcessing: ", link)

		parsedLink, _ := url.Parse(link)
		if parsedLink.Host != cfg.baseURL.Host {
			cfg.addExternalLink(parsedLink.String())
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
	if !ok {
		cfg.pages[normalizedURL] = 1
		return true
	}

	cfg.pages[normalizedURL] += 1
	return false
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
	pagesSorted := links.SortLinks(cfg.pages)
	externalLinksSorted := links.SortLinks(cfg.externalLinks)
	// Print the sorted slice                                    a
	fmt.Println("======== internal links =======")
	for _, item := range pagesSorted {
		fmt.Printf("%v link to %v%v\n", item.Visits, cfg.baseURL.String(), item.Url)
	}

	fmt.Println("======= external links ========")
	for _, item := range externalLinksSorted {
		fmt.Printf("%v link to %v\n", item.Visits, item.Url)
	}
}

func (cfg config) addExternalLink(link string) {
	_, ok := cfg.externalLinks[link]
	if !ok {
		cfg.externalLinks[link] = 1
		return
	}

	cfg.externalLinks[link]++
}

func (cfg config) outputCSV() {
	filename := "links.csv"
	writer, file, err := createCSVWriter(filename)
	if err != nil {
		fmt.Println("error creating CSV writer:", err)
		return
	}

	defer file.Close()

	header := []string{"type", "count", "url"}
	writeCSVRecord(writer, header)

	internal := links.SortLinks(cfg.pages)
	external := links.SortLinks(cfg.externalLinks)

	var records [][]string
	for _, item := range internal {
		records = append(records, []string{"internal", fmt.Sprint(item.Visits), item.Url})
	}

	for _, item := range external {
		records = append(records, []string{"external", fmt.Sprint(item.Visits), item.Url})
	}

	for _, record := range records {
		writeCSVRecord(writer, record)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("error flushing CSV writer:", err)
	}
}
