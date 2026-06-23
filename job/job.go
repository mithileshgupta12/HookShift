package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type JobStatus int

const (
	JobPending JobStatus = iota
	JobProcessing
	JobDelivered
	JobFailed
	JobDead
)

var JobStatusName = map[JobStatus]string{
	JobPending:    "pending",
	JobProcessing: "processing",
	JobDelivered:  "delivered",
	JobFailed:     "failed",
	JobDead:       "dead",
}

func (js JobStatus) String() string {
	return JobStatusName[js]
}

type Job struct {
	JobID           string
	DestinationURL  string
	Payload         json.RawMessage
	AttemptCount    uint64
	NextAttemptTime time.Time
	Status          JobStatus
}

func (j *Job) ProcessJob() error {
	client := http.Client{Timeout: time.Second * 10}

	resp, err := client.Post(j.DestinationURL, "application/json", bytes.NewBuffer(j.Payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status code %d received from destination", resp.StatusCode)
	}

	return nil
}
