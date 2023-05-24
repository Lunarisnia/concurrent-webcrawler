package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

// TODO: make this crawler into a proper one
func crawl() error {
	response, err := http.Get("https://www.google.com/")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	z := html.NewTokenizer(response.Body)

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return z.Err()
		}
		fmt.Println(z.Token())
	}
}

func main() {
	crawl()
}
