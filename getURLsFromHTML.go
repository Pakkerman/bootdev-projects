package main

import (
	"fmt"
	"net/url"
	"path"
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
		if n.Type != html.ElementNode || n.Data != "a" {
			// Continue traversing child nodes
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				traverse(c)
			}
		}
		// append url
		for _, attr := range n.Attr {
			if attr.Key != "href" {
				continue
			}

			// Resolve relative URLs
			link, err := baseURL.Parse(attr.Val)
			if err != nil {
				continue
			}
			// omit link to headings
			if link.Fragment != "" {
				continue
			}

			// omit .css, .png... etc
			ext := path.Ext(link.Path)
			if ext != "" {
				continue
			}

			urls = append(urls, link.String())
		}
	}

	// Start traversing from the root node
	traverse(doc)
	return urls, nil
}
