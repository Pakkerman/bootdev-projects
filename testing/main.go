package main

import (
	"fmt"
	"net/url"
)

func main() {
	list := []string{
		"https://blog.boot.dev/misc/building-an-nlp-engine-is-hard-but-not-as-hard-as-defining-terms/#promotion-and-solicitation-detectio",
		"https://blog.boot.dev/misc/main.css",
		"https://blog.boot.dev/about/thumbnail.png",
	}

	for _, item := range list {

		parsed, _ := url.Parse(item)
		fmt.Println("Original:\t", item)
		fmt.Println("Host:\t", parsed.Host)
		fmt.Println("Hostname:\t", parsed.Hostname())
		fmt.Println("Path:\t", parsed.Path)
		fmt.Println("RawPath:\t", parsed.RawPath)
	}
}
