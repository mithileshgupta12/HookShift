package queue

import (
	"context"

	"github.com/mithileshgupta12/hook-shift/job"
)

type Queue interface {
	Enqueue(*job.Job) error
	Dequeue(context.Context) *job.Job
	Ack(string)
	Nack(*job.Job)
}
