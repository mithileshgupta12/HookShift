package worker

import (
	"context"
	"sync"

	"github.com/mithileshgupta12/hook-shift/queue"
)

func StartPool(q queue.Queue, workerCount int, ctx context.Context, wg *sync.WaitGroup) {
	for range workerCount {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					workerJob := q.Dequeue(ctx)
					err := workerJob.ProcessJob(ctx)
					if err != nil {
						q.Nack(workerJob)
						continue
					}
					q.Ack(workerJob.JobID)
				}
			}
		})
	}
}
