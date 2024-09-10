package main

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

func getHTML(rawURL string) (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Get(rawURL)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", errors.New(res.Status)
	}

	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", errors.New("content type mismatch")
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
