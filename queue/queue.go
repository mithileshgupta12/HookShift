package queue

import "github.com/mithileshgupta12/hook-shift/job"

type Queue interface {
	Enqueue(*job.Job)
	Dequeue() <-chan *job.Job
}
