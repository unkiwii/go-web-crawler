package workerpool_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/matryer/is"

	"crawler/internal/workerpool"
)

func TestPool_everyWorkIsFinishedAfterWaiting(t *testing.T) {
	workerCount := 100

	pool := workerpool.New(workerCount / 3)
	defer pool.Stop()

	var finished uint64

	work := func(n int) func(context.Context) {
		return func(context.Context) {
			time.Sleep(100 * time.Millisecond)
			atomic.AddUint64(&finished, 1)
		}
	}

	for i := 0; i < 100; i++ {
		pool.Run(work(i))
	}

	pool.Wait()

	is.New(t).Equal(finished, uint64(workerCount)) // all work must be done by now
}

func TestPool_whenPoolIsStoppedSomeWorkCanBeLost(t *testing.T) {
	workerCount := 100

	pool := workerpool.New(workerCount / 10)

	var finished uint64

	work := func(n int) func(context.Context) {
		return func(context.Context) {
			time.Sleep(100 * time.Millisecond)
			atomic.AddUint64(&finished, 1)
		}
	}

	for i := 0; i < 100; i++ {
		pool.Run(work(i))
	}

	pool.Stop()

	is.New(t).True(finished < uint64(workerCount)) // some work must be lost
}
