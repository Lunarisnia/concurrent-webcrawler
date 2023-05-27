package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var URLPool []string

func checkIfValidURL(URI string) bool {
	return strings.Contains(URI, "https://en.")
}

func crawl(URL string, target string) (bool, error) {
	response, err := http.Get(URL)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	z := html.NewTokenizer(response.Body)

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return false, z.Err()
		}
		name, _ := z.TagName()
		if string(name) == "a" {
			for {
				key, val, more := z.TagAttr()
				if string(key) == "href" {
					if string(val) == target {
						return true, nil
					}
					if checkIfValidURL(string(val)) {
						URLPool = append(URLPool, string(val))
					}
				}

				if !more {
					break
				}
			}
		}
	}
}

func main() {
	start := time.Now()
	target := "https://en.wikipedia.org/wiki/Indonesia/"
	found, _ := crawl("https://en.wikipedia.org/wiki/Main_Page", target)
	if !found {
		for i := 0; i < len(URLPool); i++ {
			found, _ = crawl(URLPool[i], target)
			if found {
				break
			}
		}
	}
	duration := time.Since(start)
	fmt.Println(target, " Found in: ", duration)
}
