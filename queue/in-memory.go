package queue

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/mithileshgupta12/hook-shift/job"
)

type InMemoryQueue struct {
	records       chan *job.Job
	activeRecords sync.Map
	failedRecords sync.Map
	nackTimes     map[uint64]time.Duration
}

func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{
		records: make(chan *job.Job, 10_000),
		nackTimes: map[uint64]time.Duration{
			0: time.Minute,
			1: time.Minute * 5,
			2: time.Minute * 30,
			3: time.Hour * 2,
			4: time.Hour * 12,
		},
	}
}

func (imq *InMemoryQueue) Enqueue(workerJob *job.Job) error {
	imqRecords := imq.records

	select {
	case imqRecords <- workerJob:
	default:
		return errors.New("queue is full")
	}

	return nil
}

func (imq *InMemoryQueue) Dequeue(ctx context.Context) *job.Job {
	select {
	case <-ctx.Done():
		return nil
	case workerJob := <-imq.records:
		workerJob.Status = job.JobProcessing
		imq.activeRecords.Store(workerJob.JobID, workerJob)
		imq.failedRecords.Delete(workerJob.JobID)
		return workerJob
	}
}

func (imq *InMemoryQueue) Ack(jobID string) {
	imq.activeRecords.Delete(jobID)
	imq.failedRecords.Delete(jobID)
}

func (imq *InMemoryQueue) Nack(workerJob *job.Job) {
	imq.activeRecords.Delete(workerJob.JobID)
	imq.failedRecords.Delete(workerJob.JobID)

	if workerJob.AttemptCount == 5 {
		return
	}

	currentAttemptCount := workerJob.AttemptCount
	workerJob.NextAttemptTime = time.Now().Add(imq.nackTimes[workerJob.AttemptCount])
	workerJob.AttemptCount += 1
	workerJob.Status = job.JobFailed

	imq.failedRecords.Store(workerJob.JobID, workerJob)

	time.AfterFunc(imq.nackTimes[currentAttemptCount], func() {
		imq.records <- workerJob
	})
}
