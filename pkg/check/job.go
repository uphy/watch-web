package check

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	StatusRunning Status = "running"
	StatusOK      Status = "ok"
	StatusError   Status = "error"
)

type (
	// JobInfo has static information defined with config file.
	JobInfo struct {
		ID    string `json:"id"`
		Link  string `json:"link,omitempty"`
		Label string `json:"label,omitempty"`
	}
	// JobStatus has dynamic information without job result value.
	JobStatus struct {
		Status Status     `json:"status"`
		Last   *time.Time `json:"last,omitempty"`
		Error  *string    `json:"error,omitempty"`
		Count  int        `json:"count"`
	}
	Job struct {
		Info   *JobInfo `json:"info"`
		source Source
		ctx    *JobContext
	}
	JobContext struct {
		Log *logrus.Entry
	}
)

func NewJob(info *JobInfo, source Source) *Job {
	return &Job{info, source, nil}
}

func (j *Job) ID() string {
	return j.Info.ID
}

func (j *Job) failed(status *JobStatus, msg string, err error) {
	errw := fmt.Errorf("%s: %w", msg, err)
	errorString := errw.Error()
	status.Error = &errorString
	status.Status = StatusError
	j.ctx.Log.WithFields(logrus.Fields{
		"msg": msg,
		"err": errw,
	}).Error("Job execution failed.")
}

func (j *Job) String() string {
	return fmt.Sprintf("Job[id=%s, label=%s, source=%s]", j.ID(), j.Info.Label, j.source)
}
