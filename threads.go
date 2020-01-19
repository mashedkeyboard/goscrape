package main

import (
	"log"
)

func crawlThread(workerNum int){
	defer wg.Done()
	workerStates[workerNum] = 0
	log.Println("Thread ", workerNum, " now running")
	// if there are no links at all on the page, the end may already have arrived
	checkEnd(workerNum)
	// lock here to make sure we can know if the linkqueue is already done or not
	linkQueue.mux.Lock()
	if !linkQueue.done {
		// we don't need this locked the rest of the time
		linkQueue.mux.Unlock()
		for link := range linkQueue.channel {
			workerStates[workerNum] = 1
			fetchUrl(link, workerNum)
			workerStates[workerNum] = 0
		}
	} else {
		// unlock down here too, if the link queue is in fact done
		linkQueue.mux.Unlock()
	}
	log.Println("Finishing crawl thread ", workerNum)
	return
}

func checkEnd(workerNum int) {
	linkQueue.mux.Lock()
	if !linkQueue.done {
		done := true
		for i, s := range workerStates {
			if i != workerNum && s == 1 {
				done = false
			}
		}

		if len(linkQueue.channel) > 0 {
			done = false
		}

		if done {
			linkQueue.done = true
			log.Println("All done, closing queue")
			close(linkQueue.channel)
		}
	}
	linkQueue.mux.Unlock()
	return
}