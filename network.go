package main

import (
	"net/http"
	"golang.org/x/net/html"
	"log"
	"strings"
)

func fetchUrl(link linkToCrawl, workerNum int) {
	var resp *http.Response
	var doc *html.Node
	var err error

	rl.Take()
	resp, err = http.Get(link.url)
	if err != nil {
		log.Println("Errored while fetching link ", link.url, " (error: ", err, ")")
		return
	}
	
	if resp.Header.Get("Content-Type") == "text/html" || strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html;") {
		doc, err = html.Parse(resp.Body)
		if err != nil {
			log.Println("Errored while parsing body of link ", link.url, " (error: ", err, ")")
			return
		}
		log.Println("Fetched ", link.url, ", parsing")
		parseNode(doc, link, workerNum)
	} else {
		log.Println("Link ", link.url, " is not HTML, ignoring")
	}

	resp.Body.Close()

	// workerNum of -1 indicates the main thread
	if workerNum != -1 {
		checkEnd(workerNum)
	}
}