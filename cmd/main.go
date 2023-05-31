package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var URLPool []string
var globalVisitedPath = make(map[string]int)

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

func startSingleThread() {
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

func asyncCrawl(URLPools []string, target string, reportingChannel *chan Report, workerID int) error {
	pathTaken := make([]string, 0)
	for i := 0; i < len(URLPools); i++ {
		URL := URLPools[i]

		fmt.Printf("Worker %v Checking: %v\n", workerID, URL)
		if globalVisitedPath[URL] == 0 {
			pathTaken = append(pathTaken, URL)
			globalVisitedPath[URL] = 1

			response, err := http.Get(URL)
			if err != nil {
				return err
			}
			defer response.Body.Close()

			z := html.NewTokenizer(response.Body)

			for {
				tt := z.Next()
				if tt == html.ErrorToken {
					break
				}
				name, _ := z.TagName()
				if string(name) == "a" {
					for {
						key, val, more := z.TagAttr()
						if string(key) == "href" {
							uri := string(val)
							if uri == target || uri == target+"/" {
								*reportingChannel <- Report{found: true, workerID: workerID, pathTaken: pathTaken}
								return nil
							}
							if checkIfValidURL(uri) {
								URLPools = append(URLPools, "https://en.wikipedia.org"+uri)
							}
						}

						if !more {
							break
						}
					}
				}
			}
		}
	}
	return nil
}

type Report struct {
	found     bool
	workerID  int
	pathTaken []string
}

func main() {
	start := time.Now()
	target := "/wiki/List_of_Dragon_Ball_Z_Kai_episodes"
	found, _ := crawl("https://en.wikipedia.org/wiki/Special:Random", target)
	reportingChannel := make(chan Report, 1)
	var reportResult Report
	if !found {
		workerCount, _ := strconv.Atoi(os.Args[1])
		for i := 0; i < workerCount; i++ {
			go func(index int, nWorker int) {
				x := (len(URLPool) - 1) / nWorker
				_ = asyncCrawl(URLPool[index*x:(index+1)*x], target, &reportingChannel, index+1)
			}(i, workerCount)
		}
	}
	reportResult = <-reportingChannel
	if reportResult.found {
		duration := time.Since(start)
		fmt.Println("Path Taken: ")
		for ind, p := range reportResult.pathTaken {
			fmt.Printf("%v. %v\n", ind+1, p)
		}
		fmt.Println(target, " Found in: ", duration, " by worker: ", reportResult.workerID)
	} else {
		fmt.Println("Failed to find the target :(")
	}
}
