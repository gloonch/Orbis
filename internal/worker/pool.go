package worker

import (
	"context"
	"sync"
)

type Job interface {
	Do(ctx context.Context) error
}

type Pool struct {
	wg      sync.WaitGroup
	jobs    chan Job
	workers int
}

func New(workers int, capacity int) *Pool {
	return &Pool{
		jobs:    make(chan Job, capacity),
		workers: workers,
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job := <-p.jobs:
					_ = job.Do(ctx)
				}
			}
		}()
	}
}

func (p *Pool) Submit(job Job) {
	p.jobs <- job
}

func (p *Pool) Stop() {
	close(p.jobs)
	p.wg.Wait()
}
