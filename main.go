package main

import (
	"fmt"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)

	} else if 2 < len(os.Args) {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	url := os.Args[1]
	fmt.Println("starting crawl of:", url)

	pages := make(map[string]int)

	counter := 1
	crawlPage(url, url, pages, &counter)
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int, counter *int) {
	fmt.Println("\ncrawling page:", rawCurrentURL)
	fmt.Println("Pages Crawled", *counter)

	parsedBase, _ := url.Parse(rawBaseURL)

	body, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	urls, _ := getURLsFromHTML(body, rawCurrentURL)
	fmt.Printf("\tfound %v urls on this page\n", len(urls))

	// for _, item := range urls {
	// 	fmt.Println(item)
	// }
	//

	for _, link := range urls {
		fmt.Println("\tProcessing: ", link)

		parsedLink, _ := url.Parse(link)
		if parsedLink.Host != parsedBase.Host {
			continue
			fmt.Printf("\t%v going out of domain, move on\n", link)

		}

		normalizedURL, err := normalizeURL(link)
		if err != nil {
			fmt.Println(err)
			continue
		}

		_, ok := pages[normalizedURL]
		if ok {
			pages[normalizedURL] += 1
			// printMap(pages)

			continue
		}

		pages[normalizedURL] = 1
		*counter++
		crawlPage(rawBaseURL, link, pages, counter)
	}
}

func printMap(input map[string]int) {
	for key, value := range input {
		fmt.Printf("%v: %v (%v entries)\n", key, value, len(input))
	}
}
