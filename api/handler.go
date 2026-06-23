package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mithileshgupta12/hook-shift/job"
	"github.com/mithileshgupta12/hook-shift/queue"
)

type Handler struct {
	q queue.Queue
}

func NewHandler(q queue.Queue) *Handler {
	return &Handler{q}
}

func (h *Handler) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (h *Handler) handleDispatches(w http.ResponseWriter, r *http.Request) {
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

	h.q.Enqueue(&job.Job{
		JobID:           id.String(),
		DestinationURL:  req.DestinationURL,
		Payload:         req.Payload,
		AttemptCount:    0,
		NextAttemptTime: time.Now(),
		Status:          job.JobPending,
	})

	w.WriteHeader(http.StatusAccepted)
}
