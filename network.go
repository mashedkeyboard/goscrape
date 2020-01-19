package main

import (
	"net/http"
	"golang.org/x/net/html"
	"log"
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

	doc, err = html.Parse(resp.Body)
	if err != nil {
		log.Println("Errored while parsing body of link ", link.url, " (error: ", err, ")")
		return
	}

	log.Println("Fetched ", link.url, ", parsing")
	parseNode(doc, link, workerNum)
	resp.Body.Close()

	// workerNum of -1 indicates the main thread
	if workerNum != -1 {
		checkEnd(workerNum)
	}
}