package check

import (
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/robfig/cron"
)

type (
	Executor struct {
		c          *cron.Cron
		Jobs       []*Job
		initialRun bool
	}
	Job struct {
		Name     string `json:"name"`
		source   Source
		actions  []Action
		Previous *string    `json:"previous,omitempty"`
		Status   Status     `json:"status"`
		Error    error      `json:"error,omitempty"`
		Last     *time.Time `json:"last,omitempty"`
		Count    int        `json:"count"`
	}
	Source interface {
		Fetch() (string, error)
		Label() string
	}
	Action interface {
		Run(result *Result) error
	}
	Status string
)

const (
	StatusRunning Status = "running"
	StatusOK      Status = "ok"
	StatusError   Status = "error"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func NewExecutor() *Executor {
	return &Executor{
		c: cron.New(),
	}
}

func (e *Executor) Job(name string) *Job {
	for _, j := range e.Jobs {
		if j.Name == name {
			return j
		}
	}
	return nil
}

func (e *Executor) AddJob(name string, schedule string, source Source, action ...Action) error {
	if schedule == "" {
		schedule = "@every 1h"
	}
	job := &Job{
		Name:     name,
		source:   source,
		actions:  action,
		Previous: nil,
		Status:   StatusOK,
		Error:    nil,
	}
	e.Jobs = append(e.Jobs, job)
	return e.c.AddFunc(schedule, func() {
		job.Check()
	})
}

func (e *Executor) Start() {
	if e.initialRun {
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
	log.Printf("Running job: %s", j.Name)
	now := time.Now()
	j.Last = &now
	j.Count++

	j.Status = StatusRunning
	current, err := j.source.Fetch()
	if err != nil {
		j.Error = err
		j.Status = StatusError
		log.Printf("Failed to fetch %s: %v", j.Name, err)
		return
	}

	if j.Previous != nil {
		result := &Result{j.Name, j.source.Label(), *j.Previous, current}
		if err := j.doActions(result); err != nil {
			j.Status = StatusError
			j.Error = err
			log.Printf("Failed to perform action %s: %v", j.Name, err)
		}
	}
	j.Previous = &current
	j.Status = StatusOK
	log.Printf("Finished job: %s", j.Name)
}

func (j *Job) TestActions() error {
	return j.doActions(&Result{
		Name:     j.Name,
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
