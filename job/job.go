package job

import (
	"encoding/json"
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
