package main

import (
	"fmt"
)

func main() {
	// if len(os.Args) < 2 {
	// 	fmt.Println("no website provided")
	// 	os.Exit(1)
	//
	// } else if 2 < len(os.Args) {
	// 	fmt.Println("too many arguments provided")
	// 	os.Exit(1)
	// }
	//
	// url := os.Args[1]
	//
	// fmt.Println("starting crawl of:", url)

	url := "https://blog.boot.dev"
	body, err := getBody(url)
	if err != nil {
		fmt.Println(err)
	}

	urls, _ := getURLsFromHTML(body, url)
	fmt.Println(urls)
}
