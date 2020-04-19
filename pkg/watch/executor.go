package watch

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/watch/store"
)

type (
	Executor struct {
		c          *cron.Cron
		Jobs       map[string]*Job
		InitialRun bool
		store      domain.Store
		actions    []domain.Action
		log        *logrus.Logger
	}
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func NewExecutor(store domain.Store, actions []domain.Action, log *logrus.Logger) *Executor {
	return &Executor{
		c:     cron.New(),
		store: store,
		Jobs:  make(map[string]*Job),
		log:   log,
	}
}

func (e *Executor) Job(id string) *Job {
	return e.Jobs[id]
}

func (e *Executor) AddJob(job *Job, schedule *string) error {
	jobLogger := e.log.WithFields(logrus.Fields{
		"id": job.ID(),
	})
	job.ctx = &domain.JobContext{jobLogger}
	e.Jobs[job.ID()] = job
	if schedule != nil {
		return e.c.AddFunc(*schedule, func() {
			e.Check(job)
		})
	}
	return nil
}

func (e *Executor) Run() {
	if e.InitialRun {
		go e.CheckAll()
	}

	e.c.Run()
}

func (e *Executor) CheckAll() {
	ch := make(chan struct{}, 1)
	wg := new(sync.WaitGroup)
	for _, job := range e.Jobs {
		wg.Add(1)
		ch <- struct{}{}
		go func(job *Job) {
			defer wg.Done()
			e.Check(job)
			<-ch
		}(job)
	}
	wg.Wait()
}

func (e *Executor) Check(job *Job) (res *domain.Result, err error) {
	job.ctx.Log.Info("Running job.")

	// Get previous job properties
	status, err := e.store.GetStatus(job.ID())
	if status == nil {
		status = new(domain.JobStatus)
	}
	defer func() {
		// store status even if got errors for fixing broken data
		if err := e.store.SetStatus(job.ID(), status); err != nil {
			job.failed(status, "failed to set job status", err)
		}
	}()
	if err != nil && err != store.ErrNotFound {
		job.failed(status, "failed to get previous job status", err)
		return nil, err
	}
	previous, err := e.store.GetValue(job.ID())
	firstCheck := false
	if err != nil {
		if err == store.ErrNotFound {
			firstCheck = true
		} else {
			job.failed(status, "failed to load previous job value", err)
			return nil, err
		}
	}

	now := time.Now()
	status.Last = &now
	status.Count++
	status.Status = domain.StatusRunning
	status.Error = nil

	// Get current value
	current, err := job.source.Fetch(job.ctx)
	if err != nil {
		job.failed(status, "failed to fetch", err)
		return
	}
	job.ctx.Log.WithFields(logrus.Fields{
		"current": fmt.Sprintf("%#v", current),
	}).Debug("Fetched job result.")
	job.ctx.Log.Info("Fetched job result.")
	currentString := current.String()
	defer func() {
		if err := e.store.SetValue(job.ID(), currentString); err != nil {
			job.failed(status, "failed to store job value", err)
		}
	}()

	// make result
	res = domain.NewResult(job.ID(), job.Info.Label, job.Info.Link, previous, currentString, current.Type())

	// Do action
	if !firstCheck {
		if err = e.DoActions(job, res); err != nil {
			job.failed(status, "failed to perform action", err)
			return
		}
	}

	// Store job status
	status.Status = domain.StatusOK
	job.ctx.Log.WithField("result", fmt.Sprintf("%#v", res)).Debug("Finished job.")
	job.ctx.Log.Info("Finished job.")
	return
}

func (e *Executor) TestActions(job *Job) error {
	return e.DoActions(job, &domain.Result{
		JobID:    job.ID(),
		Label:    "test action",
		Link:     "https://google.com",
		Previous: "This is test action.",
		Current:  "This is test action.\naaaa",
	})
}

func (e *Executor) DoActions(job *Job, result *domain.Result) error {
	var errs error
	for _, action := range e.actions {
		if err := action.Run(job.ctx, result); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

func (e *Executor) GetJobStatus(jobID string) (*domain.JobStatus, error) {
	status, err := e.store.GetStatus(jobID)
	if err != nil {
		if err == store.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return status, nil
}

func (e *Executor) GetJobValue(jobID string) (*string, error) {
	status, err := e.store.GetValue(jobID)
	if err != nil {
		if err == store.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &status, nil
}
