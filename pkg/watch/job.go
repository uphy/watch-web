package watch

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/uphy/watch-web/pkg/domain"
)

type (
	Job struct {
		Info    *domain.JobInfo `json:"info"`
		source  domain.Source
		ctx     *domain.JobContext
		actions []domain.Action
	}
)

func NewJob(info *domain.JobInfo, source domain.Source, actions []domain.Action) *Job {
	return &Job{info, source, nil, actions}
}

func (j *Job) ID() string {
	return j.Info.ID
}

func (j *Job) failed(status *domain.JobStatus, msg string, err error) {
	errw := fmt.Errorf("%s: %w", msg, err)
	errorString := errw.Error()
	status.Error = &errorString
	status.Status = domain.StatusError
	j.ctx.Log.WithFields(logrus.Fields{
		"msg": msg,
		"err": errw,
	}).Error("Job execution failed.")
}

func (j *Job) String() string {
	return fmt.Sprintf("Job[id=%s, label=%s, source=%s]", j.ID(), j.Info.Label, j.source)
}
