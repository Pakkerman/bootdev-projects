package main

import (
	"fmt"
)

func main() {
	ch := make(chan struct{}, 5)

	for i := 0; i < 100; i++ {
		fmt.Println(i, len(ch))
		ch <- struct{}{}

		if len(ch) == 4 {
			<-ch
		}
	}
}
