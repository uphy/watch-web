package config

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/check"
	"github.com/uphy/watch-web/pkg/config/template"
)

func (c *JobConfig) addTo(ctx *template.TemplateContext, e *check.Executor) ([]*check.Job, error) {
	jobs := make([]*check.Job, 0)
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

func evaluateItemAsTemplate(ctx *template.TemplateContext, v interface{}) (interface{}, error) {
	m, ok := v.(map[string]interface{})
	if ok {
		evaluated := make(map[string]interface{})
		for key, value := range m {
			ekey, err := template.TemplateString(key).Evaluate(ctx)
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
	e, err := template.TemplateString(fmt.Sprint(v)).Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (c *JobConfig) addOne(ctx *template.TemplateContext, e *check.Executor) (*check.Job, error) {
	source, err := c.Source.Source(ctx)
	if err != nil {
		return nil, err
	}
	actions := []check.Action{}
	for _, actionConfig := range c.Actions {
		action, err := actionConfig.Action(ctx)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
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
	job, err := e.AddJob(id, schedule, label, link, source, actions...)
	if err != nil {
		return nil, err
	}
	return job, nil
}
