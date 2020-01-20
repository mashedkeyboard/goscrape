package main

import (
	"go.uber.org/ratelimit"
	"net/url"
	"sync"
)

const maxDepth, queueLength, threads int = 2, 4096, 5

var linkQueue syncedLinkQueue
var wg sync.WaitGroup
var startUrl url.URL
var urls urlListContainer
var visited visitLog
var workerStates []int
var rl ratelimit.Limiter
