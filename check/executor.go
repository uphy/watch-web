package check

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/robfig/cron"
)

type (
	Executor struct {
		c          *cron.Cron
		Jobs       map[string]*Job
		InitialRun bool
		store      Store
	}
	Source interface {
		Fetch() (string, error)
	}
	Action interface {
		Run(result *Result) error
	}
	Status string
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func NewExecutor(store Store) *Executor {
	return &Executor{
		c:     cron.New(),
		store: store,
		Jobs:  make(map[string]*Job),
	}
}

func (e *Executor) Job(id string) *Job {
	return e.Jobs[id]
}

func (e *Executor) AddJob(id, schedule, label, link string, source Source, action ...Action) error {
	if schedule == "" {
		schedule = "@every 1h"
	}
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
	}
	job.RestoreState()
	e.Jobs[id] = job
	return e.c.AddFunc(schedule, func() {
		job.Check()
	})
}

func (e *Executor) Run() {
	if e.InitialRun {
		go e.checkAll()
	}

	e.c.Run()
}

func (e *Executor) checkAll() {
	ch := make(chan struct{}, 3)
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
	return nil
}

func (j *Job) StoreState() error {
	if err := j.store.SetJob(j.ID, j); err != nil {
		j.failed("failed to store current value", err)
		return err
	}
	return nil
}

func (j *Job) Check() (result *Result) {
	log.Printf("Running job: %s", j.ID)
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
	current, err := j.source.Fetch()
	if err != nil {
		j.failed("failed to fetch", err)
		return
	}
	current = strings.Trim(current, " \t\n")

	// make result
	var previous string
	if j.Previous != nil {
		previous = *j.Previous
	} else {
		previous = ""
	}
	result = &Result{j.ID, j.Label, j.Link, previous, current}

	// Do action
	if j.Previous != nil {
		if err := j.doActions(result); err != nil {
			j.failed("failed to perform action", err)
		}
	}

	// Store job status
	j.Previous = &current
	j.Status = StatusOK
	j.StoreState()
	log.Printf("Finished job: id=%s, result=%v", j.ID, result)
	return
}

func (j *Job) TestActions() error {
	return j.doActions(&Result{
		JobID:    j.ID,
		Label:    "test action",
		Link:     "https://google.com",
		Previous: "This is test action.",
		Current:  "This is test action.\naaaa",
	})
}

func (j *Job) doActions(result *Result) error {
	var errs error
	for _, action := range j.actions {
		if err := action.Run(result); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}
