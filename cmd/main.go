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
	hasWiki := strings.Contains(URI, "/wiki")
	if !hasWiki {
		return false
	}

	return URI[:5] == "/wiki"
}

func crawl(URL string, target string) (bool, error) {
	fmt.Println("Checking: ", URL)
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
					uri := string(val)
					if uri == target || uri == target+"/" {
						return true, nil
					}
					if checkIfValidURL(uri) {
						URLPool = append(URLPool, "https://en.wikipedia.org"+uri)
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
	target := "/wiki/Indonesia"
	found, _ := crawl("https://en.wikipedia.org/wiki/Pro_Football_Hall_of_Fame", target)
	if !found {
		for i := 0; i < len(URLPool); i++ {
			found, _ = crawl(URLPool[i], target)
			if found {
				break
			}
		}
	}

	if found {
		duration := time.Since(start)
		fmt.Println(target, " Found in: ", duration)
	} else {
		fmt.Println("Failed to find the target :(")
	}
}
