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
		initialRun bool
		store      Store
	}
	Job struct {
		Name     string `json:"name"`
		source   Source
		actions  []Action
		Link     string     `json:"link,omitempty"`
		Label    string     `json:"label,omitempty"`
		Previous *string    `json:"previous,omitempty"`
		Status   Status     `json:"status"`
		Error    error      `json:"error,omitempty"`
		Last     *time.Time `json:"last,omitempty"`
		Count    int        `json:"count"`
		store    Store
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

func NewExecutor(store Store) *Executor {
	return &Executor{
		c:     cron.New(),
		store: store,
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

func (e *Executor) AddJob(name, schedule, label, link string, source Source, action ...Action) error {
	if schedule == "" {
		schedule = "@every 1h"
	}
	job := &Job{
		Name:     name,
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
	current = strings.Trim(current, " \t\n")
	prev, err := j.store.Get(j.Name)
	if err != nil {
		log.Println("failed to get previous value: ", err)
		return
	}

	if prev != nil {
		result := &Result{j.Name, j.Label, j.Link, *prev, current}
		if err := j.doActions(result); err != nil {
			j.Status = StatusError
			j.Error = err
			log.Printf("Failed to perform action %s: %v", j.Name, err)
		}
	}
	if err := j.store.Set(j.Name, current); err != nil {
		log.Println("failed to set current value: ", err)
		return
	}
	j.Previous = &current
	j.Status = StatusOK
	log.Printf("Finished job: %s", j.Name)
}

func (j *Job) TestActions() error {
	return j.doActions(&Result{
		Name:     j.Name,
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
