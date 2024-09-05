package main

import (
	"net/url"
	"strings"
)

func normalizeURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	normalizedURL := parsedURL.Host + strings.TrimRight(parsedURL.Path, "/")
	return normalizedURL, nil
}
