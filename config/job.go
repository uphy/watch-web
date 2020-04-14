package config

import (
	"fmt"

	"github.com/uphy/watch-web/check"
)

func (c *JobConfig) addTo(ctx *TemplateContext, e *check.Executor) error {
	if len(c.WithItems) == 0 {
		return c.addOne(ctx, e)
	}
	for itemIndex, item := range c.WithItems {
		evaluatedItem, err := evaluateItemAsTemplate(ctx, item)
		if err != nil {
			return err
		}
		ctx.PushScope()
		ctx.Set("itemIndex", itemIndex)
		ctx.Set("item", evaluatedItem)
		c.addOne(ctx, e)
		ctx.PopScope()
	}
	return nil
}

func evaluateItemAsTemplate(ctx *TemplateContext, v interface{}) (interface{}, error) {
	m, ok := v.(map[string]interface{})
	if ok {
		evaluated := make(map[string]interface{})
		for key, value := range m {
			ekey, err := TemplateString(key).Evaluate(ctx)
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
	e, err := TemplateString(fmt.Sprint(v)).Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (c *JobConfig) addOne(ctx *TemplateContext, e *check.Executor) error {
	source, err := c.Source.Source(ctx)
	if err != nil {
		return err
	}
	actions := []check.Action{}
	for _, actionConfig := range c.Actions {
		action, err := actionConfig.Action(ctx)
		if err != nil {
			return err
		}
		actions = append(actions, action)
	}
	id, err := c.ID.Evaluate(ctx)
	if err != nil {
		return err
	}
	schedule, err := c.Schedule.Evaluate(ctx)
	if err != nil {
		return err
	}
	label, err := c.Label.Evaluate(ctx)
	if err != nil {
		return err
	}
	link, err := c.Link.Evaluate(ctx)
	if err != nil {
		return err
	}
	if err := e.AddJob(id, schedule, label, link, source, actions...); err != nil {
		return err
	}
	return nil
}
