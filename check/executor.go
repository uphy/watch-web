package check

import (
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/robfig/cron"
)

type (
	Executor struct {
		c          *cron.Cron
		Jobs       []*Job
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
	}
}

func (e *Executor) Job(id string) *Job {
	for _, j := range e.Jobs {
		if j.ID == id {
			return j
		}
	}
	return nil
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
	e.Jobs = append(e.Jobs, job)
	log.Printf("Job added: id=%s, label=%s", id, label)
	return e.c.AddFunc(schedule, func() {
		job.Check()
	})
}

func (e *Executor) Start() {
	if e.InitialRun {
		go e.checkAll()
	}

	e.c.Start()
}

func (e *Executor) checkAll() {
	for _, job := range e.Jobs {
		go job.Check()
	}
}

func (j *Job) Check() {
	log.Printf("Running job: %s", j.ID)
	// Restore job status
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

	// Do action
	if j.Previous != nil {
		result := &Result{j.ID, j.Label, j.Link, *j.Previous, current}
		if err := j.doActions(result); err != nil {
			j.failed("failed to perform action", err)
		}
	}

	// Store job status
	j.Previous = &current
	j.Status = StatusOK
	if err := j.store.SetJob(j.ID, j); err != nil {
		j.failed("failed to store current value", err)
		return
	}
	log.Printf("Finished job: %s", j.ID)
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
