package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	reader := strings.NewReader(htmlBody)
	doc, err := html.Parse(reader)
	if err != nil {
		fmt.Println(err)
	}

	var urls []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		// Check if the node is an <a> element with an href attribute
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					// Resolve relative URLs
					link, err := baseURL.Parse(attr.Val)
					if err == nil {
						urls = append(urls, link.String())
					}
				}
			}
		}
		// Continue traversing child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	// Start traversing from the root node
	traverse(doc)
	return urls, nil
}
