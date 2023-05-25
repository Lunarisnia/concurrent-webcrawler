package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func crawl() error {
	response, err := http.Get("http://localhost:3000")
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
		name, _ := z.TagName()
		if string(name) == "a" {
			for {
				key, val, more := z.TagAttr()
				if string(key) == "href" {
					fmt.Println(string(val))
				}

				if !more {
					break
				}
			}
		}
	}
}

func main() {
	crawl()
}
