package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mithileshgupta12/hook-shift/job"
	"github.com/mithileshgupta12/hook-shift/queue"
)

func main() {
	portAddr := flag.Int("port", 9000, "port defines the port on which the api runs")
	workersAddr := flag.Int("workers", 5, "workers defines the number of worker nodes for processing the jobs")
	flag.Parse()

	inMemoryQueue := queue.NewInMemoryQueue()

	for range *workersAddr {
		go func() {
			for {
				workerJob := inMemoryQueue.Dequeue()
				err := workerJob.ProcessJob()
				if err != nil {
					inMemoryQueue.Nack(workerJob)
					continue
				}
				inMemoryQueue.Ack(workerJob.JobID)
			}
		}()
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /v1/dispatches", func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			DestinationURL string          `json:"destination_url"`
			Payload        json.RawMessage `json:"payload"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("some error occurred while decoding the request body: %v", err)
			http.Error(w, "some error occurred while decoding the request body", http.StatusBadRequest)
			return
		}

		id := uuid.New()

		inMemoryQueue.Enqueue(&job.Job{
			JobID:           id.String(),
			DestinationURL:  req.DestinationURL,
			Payload:         req.Payload,
			AttemptCount:    0,
			NextAttemptTime: time.Now(),
			Status:          job.JobPending,
		})

		w.WriteHeader(http.StatusAccepted)
	})

	log.Printf("listening on port %d", *portAddr)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *portAddr), mux))
}
