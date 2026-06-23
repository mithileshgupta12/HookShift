package queue

import (
	"sync"

	"github.com/mithileshgupta12/hook-shift/job"
)

type InMemoryQueue struct {
	mu      sync.Mutex
	records chan *job.Job
}

func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{
		records: make(chan *job.Job, 10_000),
	}
}

func (imq *InMemoryQueue) Enqueue(job *job.Job) {
	imq.records <- job
}

func (imq *InMemoryQueue) Dequeue() <-chan *job.Job {
	return imq.records
}
