package main

import (
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strings"
)

func fetchUrl(link linkToCrawl, workerNum int) {
	var resp *http.Response
	var doc *html.Node
	var err error

	rl.Take()
	resp, err = http.Get(link.url)
	if err == nil {
		if resp.Header.Get("Content-Type") == "text/html" || strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html;") {
			doc, err = html.Parse(resp.Body)
			if err == nil {
				log.Println("Fetched ", link.url, ", parsing")
				parseNode(doc, link, workerNum)
			} else {
				log.Println("Errored while parsing body of link ", link.url, " (error: ", err, ")")
			}
		} else {
			log.Println("Link ", link.url, " is not HTML, ignoring")
		}
	
		resp.Body.Close()
	} else {
		log.Println("Errored while fetching link ", link.url, " (error: ", err, ")")
	}

	// workerNum of -1 indicates the main thread
	if workerNum != -1 {
		checkEnd(workerNum)
	}

	return
}
