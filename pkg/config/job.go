package config

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/watch"
)

type (
	JobConfig struct {
		ID        domain.TemplateString `json:"id"`
		Label     domain.TemplateString `json:"label"`
		Link      domain.TemplateString `json:"link"`
		Source    *SourceConfig         `json:"source,omitempty"`
		Schedule  domain.TemplateString `json:"schedule,omitempty"`
		WithItems []interface{}         `json:"with_items,omitempty"`
	}
)

func (c *JobConfig) addTo(ctx *domain.TemplateContext, e *watch.Executor) ([]*watch.Job, error) {
	jobs := make([]*watch.Job, 0)
	if len(c.WithItems) == 0 {
		job, err := c.addOne(ctx, e)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	} else {
		for itemIndex, item := range c.WithItems {
			evaluatedItem, err := evaluateItemAsTemplate(ctx, item)
			if err != nil {
				return nil, err
			}
			ctx.PushScope()
			ctx.Set("itemIndex", itemIndex)
			ctx.Set("item", evaluatedItem)
			j, err := c.addOne(ctx, e)
			if err != nil {
				return nil, err
			}
			jobs = append(jobs, j)
			ctx.PopScope()
		}
	}
	return jobs, nil
}

func evaluateItemAsTemplate(ctx *domain.TemplateContext, v interface{}) (interface{}, error) {
	m, ok := v.(map[string]interface{})
	if ok {
		evaluated := make(map[string]interface{})
		for key, value := range m {
			ekey, err := domain.TemplateString(key).Evaluate(ctx)
			if err != nil {
				return nil, err
			}
			evalue, err := evaluateItemAsTemplate(ctx, value)
			if err != nil {
				return nil, err
			}
			evaluated[ekey] = evalue
		}
		return evaluated, nil
	}
	e, err := domain.TemplateString(fmt.Sprint(v)).Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (c *JobConfig) addOne(ctx *domain.TemplateContext, e *watch.Executor) (*watch.Job, error) {
	source, err := c.Source.Source(ctx)
	if err != nil {
		return nil, err
	}
	id, err := c.ID.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	schedule, err := c.Schedule.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	label, err := c.Label.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	link, err := c.Link.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	job := watch.NewJob(&domain.JobInfo{
		ID:    id,
		Label: label,
		Link:  link,
	}, source)

	if err := e.AddJob(job, &schedule); err != nil {
		return nil, err
	}
	return job, nil
}
