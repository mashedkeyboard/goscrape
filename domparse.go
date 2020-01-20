package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/url"
	"strings"
)

func parseNode(n *html.Node, link linkToCrawl, workerNum int) {
	// parse this node itself if it's an anchor
	if n.Type == html.ElementNode && n.Data == "a" && link.depth+1 <= maxDepth {
		for _, a := range n.Attr {
			// only spider links that are valid and not just anchors
			if a.Key == "href" && !strings.HasPrefix(a.Val, "#") {
				// try to parse the URL
				urlval := a.Val
				url, err := url.Parse(urlval)
				if err == nil {
					if url.Hostname() != startUrl.Hostname() {
						// if it's not got the same hostname, it can either be a link elsewhere or a relative path
						if len(url.Scheme) == 0 {
							// if it doesn't validate as a request URL, but did validate as a URL,
							// and has no scheme, it's probably a relative path
							urlval = fmt.Sprint(startUrl.Scheme, "://", startUrl.Host, urlval)
						} else {
							// it's probably a link elsewhere or a non-web scheme
							continue
						}
					} else {
						urlval = url.String()
					}
					visited.mux.Lock()
					if _, exists := visited.urls[urlval]; !exists {
						visited.urls[urlval] = true
						visited.mux.Unlock()

						x := *new(linkToCrawl)
						x.url = urlval
						x.parentLink = &link
						x.depth = link.depth + 1
						urls.addLink(x)
						linkQueue.channel <- x
						// log.Println("Added ", x.url, " to link queue on worker num ", workerNum)
					} else {
						visited.mux.Unlock()
					}
				}
			}
		}
	}
	// parse all the children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseNode(c, link, workerNum)
	}
}
