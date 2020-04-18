package check

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/uphy/watch-web/pkg/result"
	"github.com/uphy/watch-web/pkg/value"
)

type (
	Executor struct {
		c          *cron.Cron
		Jobs       map[string]*Job
		InitialRun bool
		store      Store
		log        *logrus.Logger
	}
	Source interface {
		Fetch(ctx *JobContext) (value.Value, error)
	}
	Action interface {
		Run(ctx *JobContext, result *result.Result) error
	}
	Status string
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func NewExecutor(store Store, log *logrus.Logger) *Executor {
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

func (e *Executor) AddJob(id, schedule, label, link string, source Source, action ...Action) (*Job, error) {
	if schedule == "" {
		schedule = "@every 1h"
	}
	jobLogger := e.log.WithFields(logrus.Fields{
		"id": id,
	})
	job := &Job{
		ID:       id,
		source:   source,
		actions:  action,
		Label:    label,
		Link:     link,
		Previous: nil,
		Status:   StatusOK,
		Error:    nil,
		store:    e.store,
		ctx:      &JobContext{jobLogger},
	}
	job.RestoreState()
	e.Jobs[id] = job
	return job, e.c.AddFunc(schedule, func() {
		job.Check()
	})
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
			job.Check()
			<-ch
		}(job)
	}
	wg.Wait()
}

func (j *Job) RestoreState() error {
	if err := j.store.GetJob(j.ID, j); err != nil {
		j.failed("failed to get previous job status", err)
		return err
	}
	j.ctx.Log.WithField("job", j).Debug("Restored job.")
	j.ctx.Log.Info("Restored job.")
	return nil
}

func (j *Job) StoreState() error {
	if err := j.store.SetJob(j.ID, j); err != nil {
		j.failed("failed to store current value", err)
		return err
	}
	j.ctx.Log.WithField("job", j).Debug("Stored job.")
	j.ctx.Log.Info("Stored job.")
	return nil
}

func (j *Job) Check() (res *result.Result) {
	j.ctx.Log.Info("Running job.")
	j.RestoreState()

	if err := j.store.GetJob(j.ID, j); err != nil {
		j.failed("failed to get previous job status", err)
	}

	now := time.Now()
	j.Last = &now
	j.Count++
	j.Status = StatusRunning
	j.Error = nil

	// Get current value
	current, err := j.source.Fetch(j.ctx)
	if err != nil {
		j.failed("failed to fetch", err)
		return
	}
	j.ctx.Log.WithFields(logrus.Fields{
		"current": fmt.Sprintf("%#v", current),
	}).Debug("Fetched job result.Result.")
	j.ctx.Log.Info("Fetched job result.Result.")

	// make result.Result
	var previous string
	if j.Previous != nil {
		previous = *j.Previous
	} else {
		previous = ""
	}
	currentString := current.String()
	res = result.New(j.ID, j.Label, j.Link, previous, currentString)

	// Do action
	if j.Previous != nil {
		if err := j.doActions(res); err != nil {
			j.failed("failed to perform action", err)
		}
	}

	// Store job status
	j.Previous = &currentString
	j.Status = StatusOK
	j.StoreState()
	j.ctx.Log.WithField("result.Result", fmt.Sprintf("%#v", res)).Debug("Finished job.")
	j.ctx.Log.Info("Finished job.")
	return
}

func (j *Job) TestActions() error {
	return j.doActions(&result.Result{
		JobID:    j.ID,
		Label:    "test action",
		Link:     "https://google.com",
		Previous: "This is test action.",
		Current:  "This is test action.\naaaa",
	})
}

func (j *Job) doActions(result *result.Result) error {
	var errs error
	for _, action := range j.actions {
		if err := action.Run(j.ctx, result); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}
