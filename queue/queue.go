package queue

import "github.com/mithileshgupta12/hook-shift/job"

type Queue interface {
	Enqueue(*job.Job) error
	Dequeue() *job.Job
	Ack(string)
	Nack(*job.Job)
}
