package main

import (
	"fmt"
	"go.uber.org/ratelimit"
	"log"
	"net/url"
	"os"
	"strings"
)

func main() {
	// initialise the visit log, queue, and rate limiter
	visited.urls = make(map[string]bool)
	linkQueue.channel = make(chan linkToCrawl, queueLength)
	rl = ratelimit.New(2) // per second
	if len(os.Args) != 2 {
		log.Fatal("Please specify a start URL as the only argument to this program.")
	}
	startLink := *(new(linkToCrawl))
	startLink.url = os.Args[1]

	u, err := url.Parse(startLink.url)
	if err != nil {
		log.Fatal("That wasn't a valid start URL. Please specify a valid start URL.")
	}
	startUrl = *u

	log.Println("Fetching first URL")
	urls.list.url = startUrl.String()
	// thread -1 indicates the main thread, so it won't try and check it or shut it down
	fetchUrl(startLink, -1)

	log.Println("Starting threads")
	workerStates = make([]int, threads)
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go crawlThread(i)
	}
	wg.Wait() // pauses the main thread here until all of the workers have returned
	fmt.Println(urls.list.url)
	printUrls(urls.list.suburls, 1)
	return
}

func printUrls(suburls map[string]*urlList, level int) {
	for url, urlList := range suburls {
		fmt.Println(strings.Repeat("-", level), url)
		if len((*urlList).suburls) > 0 {
			printUrls((*urlList).suburls, level+1)
		}
	}
}
