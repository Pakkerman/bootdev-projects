package main

import (
	"fmt"
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

	body, err := getHTML(url)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(body)
	urls, _ := getURLsFromHTML(body, url)
	for _, url := range urls {
		fmt.Println(url)
	}
}
