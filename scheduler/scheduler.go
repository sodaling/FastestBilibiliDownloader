package scheduler

import (
	"context"
	"simple-golang-crawler/engine"
)

type ConcurrentScheduler struct {
	RequestsChan chan *engine.Request
	WorkerChan   chan chan *engine.Request
}

func NewConcurrentScheduler() engine.Scheduler {
	return &ConcurrentScheduler{}
}

func (s *ConcurrentScheduler) Run(ctx context.Context) {
	s.WorkerChan = make(chan chan *engine.Request)
	s.RequestsChan = make(chan *engine.Request)
	go func() {
		var workerQ []chan *engine.Request
		var requestQ []*engine.Request
	loop:
		for {
			var readyWorker chan *engine.Request
			var readyRequest *engine.Request
			if len(workerQ) > 0 && len(requestQ) > 0 {
				readyWorker = workerQ[0]
				readyRequest = requestQ[0]
			}
			select {
			case readyRequest = <-s.RequestsChan:
				requestQ = append(requestQ, readyRequest)
			case readyWorker = <-s.WorkerChan:
				workerQ = append(workerQ, readyWorker)
			case readyWorker <- readyRequest:
				requestQ = requestQ[1:]
				workerQ = workerQ[1:]
			case <-ctx.Done():
				break loop
			}
		}
	}()
}

func (s *ConcurrentScheduler) GetWorkerChan() chan *engine.Request {
	return make(chan *engine.Request)
}

func (s *ConcurrentScheduler) Submit(req *engine.Request) {
	s.RequestsChan <- req
}

func (s *ConcurrentScheduler) Ready(worker chan *engine.Request) {
	s.WorkerChan <- worker
}
