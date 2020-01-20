package main

import (
	"log"
	"sync"
)

// Contains our urlLists and lets us lock them using sync's mutex lock.
// This means we can prevent two threads trying to write to a urlList at once.
type urlListContainer struct {
	list urlList
	mux  sync.Mutex
}

// The actual urlList. These have a url, which would be the root URL initially,
// and then contain a map of other urlLists, whose keys are the URLs that they have.
type urlList struct {
	url     string
	suburls map[string]*urlList
}

// linkToCrawl is - shockingly enough - a link that we're going to crawl.
// It has a URL, a parent so it can be ordered, and a depth so we know when to stop.
type linkToCrawl struct {
	url        string
	parentLink *linkToCrawl
	depth      int
}

// A syncedLinkQueue is a struct of a channel (effectively used as a queue for links),
// a variable to let us know if we're done using the queue (and thereby done crawling),
// and a mutex lock so we can lock on a single thread checking whether it's done
// or closing the channel.
type syncedLinkQueue struct {
	channel chan linkToCrawl
	done    bool
	mux     sync.Mutex
}

// A visitLog is a log of visited URLs. It has a map of URL strings against booleans.
// The boolean here will always be true for anything that exists in the map. It's
// only a map because Go has no native indexed list other than a map.
// The mutex is here to ensure that we don't have a race condition of a thread
// trying to put a visit into the log and a thread trying to check.
type visitLog struct {
	urls map[string]bool
	mux  sync.Mutex
}

// Acts on a urlListContainer to add a link to the underlying urlList
// On the container rather than the list itself to allow for the mutex lock
func (u *urlListContainer) addLink(link linkToCrawl) {
	var linkTree []linkToCrawl = []linkToCrawl{link}
	var workingLink linkToCrawl = link
	for {
		if workingLink.parentLink != nil {
			linkTree = append(linkTree, *(workingLink.parentLink))
			workingLink = *(workingLink.parentLink)
		} else {
			break
		}
	}
	if linkTree[len(linkTree)-1].url != startUrl.String() {
		// something is very wrong, this should never ever occur
		// if this happens - welp, we need to quit
		log.Fatal("linkTree initial URL didn't match start URL for link ", link.url)
	} else {
		u.mux.Lock()
		// workingUrlList is a valid copy of the entire list.
		var workingUrlList urlList = u.list
		// currentUrlList is a pointer to where we're working in the list.
		var currentUrlList *urlList = &workingUrlList
		var ok bool
		var suburl *urlList
		// going downwards from the top of the linkTree
		// we go from len(linkTree) - 2 because the top url isn't a suburl at all
		for i := len(linkTree) - 2; i >= 0; i-- {
			if (*currentUrlList).suburls == nil {
				(*currentUrlList).suburls = make(map[string]*urlList)
				ok = false
			} else {
				// get the suburl for this url, if it exists, and set the
				// pointer for the current list to it
				//
				// if there's no suburls, as above, then this can't happen
				suburl, ok = (*currentUrlList).suburls[linkTree[i].url]
			}
			if ok {
				// set the currentUrlList to the pointer to the suburl we got earlier
				currentUrlList = suburl
			} else {
				// if not, create it
				sublist := new(urlList)
				sublist.url = linkTree[i].url
				// add the new urlList to the existing part of the workingUrlList through the pointer
				(*currentUrlList).suburls[linkTree[i].url] = sublist
				// then reset the currentUrlList pointer to the current head
				currentUrlList = (*currentUrlList).suburls[linkTree[i].url]
			}
		}
		// now the link should have been added, so setup u again!
		u.list = workingUrlList
		u.mux.Unlock()
	}
	return
}
