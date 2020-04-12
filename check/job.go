package check

import (
	"fmt"
	"log"
	"time"
)

const (
	StatusRunning Status = "running"
	StatusOK      Status = "ok"
	StatusError   Status = "error"
)

type (
	// Job struct represents the scheduled job.
	// Public fields should be with `json` tag because this struct is also used as DTO.
	Job struct {
		ID       string `json:"id"`
		source   Source
		actions  []Action
		Link     string     `json:"link,omitempty"`
		Label    string     `json:"label,omitempty"`
		Previous *string    `json:"previous,omitempty"`
		Status   Status     `json:"status"`
		Error    *string    `json:"error,omitempty"`
		Last     *time.Time `json:"last,omitempty"`
		Count    int        `json:"count"`
		store    Store
	}
)

func (j *Job) failed(msg string, err error) {
	errw := fmt.Errorf("%s: %w", msg, err)
	errorString := errw.Error()
	j.Error = &errorString
	j.Status = StatusError
	log.Printf("%s: id=%s, err=%v", msg, j.ID, errw)
}
