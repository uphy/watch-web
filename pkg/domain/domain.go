package domain

import (
	"time"

	"github.com/uphy/watch-web/pkg/domain/value"

	"github.com/sirupsen/logrus"
)

const (
	StatusRunning Status = "running"
	StatusOK      Status = "ok"
	StatusError   Status = "error"
)

type (
	Status string
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
	JobContext struct {
		Log *logrus.Entry
	}
	Source interface {
		Fetch(ctx *JobContext) (value.Value, error)
	}
	Transformer interface {
		Transform(ctx *JobContext, v value.Value) (value.Value, error)
	}
	Action interface {
		Run(ctx *JobContext, result *Result) error
	}
	Store interface {
		GetJobValue(jobID string) (string, error)
		SetJobValue(jobID string, value string) error
		GetJobStatus(jobID string) (*JobStatus, error)
		SetJobStatus(jobID string, status *JobStatus) error
		SetTemp(key string, value string, expire time.Duration) error
		Get(key string) (string, error)
	}
	ScriptEngine interface {
		NewScript(script string) (Script, error)
	}
	Script interface {
		Evaluate(args map[string]interface{}) (interface{}, error)
	}
)

func NewDefaultJobContext() *JobContext {
	return &JobContext{logrus.NewEntry(logrus.New())}
}
