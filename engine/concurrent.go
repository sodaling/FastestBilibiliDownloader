package engine

import "context"

type ConcurrentEngine struct {
	WorkerCount int
	Scheduler   Scheduler
	ItemChan    chan Item
}

func NewConcurrentEngine(workerCount int, scheduler Scheduler, itemChan chan Item) *ConcurrentEngine {
	return &ConcurrentEngine{WorkerCount: workerCount, Scheduler: scheduler, ItemChan: itemChan}
}

type Scheduler interface {
	Run(context.Context)
	GetWorkerChan() chan *Request
	Submit(*Request)
	WorkerReadyNotifier
}

type WorkerReadyNotifier interface {
	Ready(chan *Request)
}

func (c *ConcurrentEngine) Run(seed ...*Request) {
	requestCount := 0
	resultChan := make(chan ParseResult)
	ctx, cancel := context.WithCancel(context.Background())
	c.Scheduler.Run(ctx)

	for i := 0; i < c.WorkerCount; i++ {
		CreateWorker(resultChan, c.Scheduler.GetWorkerChan(), c.Scheduler)
	}

	for _, req := range seed {
		hasVisited(req.Url)
		requestCount += 1
		c.Scheduler.Submit(req)
	}

	for {
		result := <-resultChan

		for _, item := range result.Items {
			go func(item Item) {
				c.ItemChan <- item
			}(item)
		}

		for _, req := range result.Requests {
			if hasVisited(req.Url) {
				continue
			} else {
				requestCount += 1
				c.Scheduler.Submit(req)
			}
		}
		requestCount -= 1
		if requestCount == 0 {
			break
		}
	}

	cancel()
}

var urlVisited = make(map[string]struct{})

func hasVisited(url string) bool {
	if _, ok := urlVisited[url]; ok {
		return true
	} else {
		urlVisited[url] = struct{}{}
	}
	return false

}

func CreateWorker(out chan ParseResult, in chan *Request, notifier WorkerReadyNotifier) {
	go func() {
		for {
			notifier.Ready(in)
			req := <-in
			ret, _ := work(req)
			out <- ret
		}
	}()
}

func work(request *Request) (ParseResult, error) {
	content, ok := request.FetchFun(request.Url)
	if ok != nil {
		return ParseResult{}, ok
	}
	result := request.ParseFunction(content, request.Url)
	return result, nil
}
