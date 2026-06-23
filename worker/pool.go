package worker

import "github.com/mithileshgupta12/hook-shift/queue"

func StartPool(q queue.Queue, workerCount int) {
	for range workerCount {
		go func() {
			for {
				workerJob := q.Dequeue()
				err := workerJob.ProcessJob()
				if err != nil {
					q.Nack(workerJob)
					continue
				}
				q.Ack(workerJob.JobID)
			}
		}()
	}
}
