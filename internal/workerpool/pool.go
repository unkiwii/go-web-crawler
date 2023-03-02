package workerpool

import (
	"context"
	"sync"
)

type Pool struct {
	workersCount int

	ctx         context.Context
	stopWorkers func()

	workQueue   chan work
	allWorkDone sync.WaitGroup
}

func New(workersCount int) *Pool {
	ctx, stopWorkers := context.WithCancel(context.Background())

	pool := Pool{
		workersCount: workersCount,
		stopWorkers:  stopWorkers,
		workQueue:    make(chan work, workersCount),
	}

	go pool.startWorkers(ctx)

	return &pool
}

func (p *Pool) Run(work func(context.Context)) {
	p.allWorkDone.Add(1)
	p.workQueue <- func(ctx context.Context) {
		work(ctx)
		p.allWorkDone.Done()
	}
}

func (p *Pool) Stop() {
	p.stopWorkers()
}

func (p *Pool) Wait() {
	p.allWorkDone.Wait()
}

func (p *Pool) startWorkers(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < p.workersCount; i++ {
		wg.Add(1)
		go p.worker(ctx, &wg, p.workQueue)
	}

	wg.Wait()
}

type work func(context.Context)

func (p *Pool) worker(ctx context.Context, wg *sync.WaitGroup, works <-chan work) {
	defer wg.Done()
	for {
		select {
		case work, ok := <-works:
			if !ok {
				return
			}
			work(ctx)
		case <-ctx.Done():
			return
		}
	}
}
