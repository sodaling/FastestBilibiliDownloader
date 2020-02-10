package engine

import "simple-golang-crawler/fetcher"

type ConcurrentEngine struct {
	WorkerCount int
	Scheduler   Scheduler
	ItemChan    chan Item
}

type Scheduler interface {
	Run()
	GetWorkerChan() chan Request
	Submit(Request)
	WorkerReadyNotifier
}

type WorkerReadyNotifier interface {
	Ready(chan Request)
}

func (c *ConcurrentEngine) Run(seed ...Request) {
	resultChan := make(chan ParseResult)
	c.Scheduler.Run()

	for i := 0; i < c.WorkerCount; i++ {
		CreateWorker(resultChan, c.Scheduler.GetWorkerChan(), c.Scheduler)
	}

	for _, req := range seed {
		c.Scheduler.Submit(req)
	}

	for {
		result := <-resultChan

		for _, item := range result.Items {
			go func() {
				c.ItemChan <- item
			}()
		}

		for _, req := range result.Requests {
			if hasVisited(req.Url) {
				continue
			} else {
				c.Scheduler.Submit(req)
			}
		}
	}

}

var urlVisited map[string]struct{}

func hasVisited(url string) bool {
	if _, ok := urlVisited[url]; ok {
		return true
	} else {
		urlVisited[url] = struct{}{}
	}
	return false

}

func CreateWorker(out chan ParseResult, in chan Request, notifier WorkerReadyNotifier) {
	go func() {
		for {
			notifier.Ready(in)
			req := <-in
			ret, err := work(req)
			if err != nil {
				continue
			}
			out <- ret
		}
	}()
}

func work(request Request) (ParseResult, error) {
	content, ok := fetcher.Fetch(request.Url)
	if ok != nil {
		return ParseResult{}, ok
	}
	result := request.ParseFunction(content, request.Url)
	return result, nil
}
