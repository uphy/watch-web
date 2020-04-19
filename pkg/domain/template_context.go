package domain

type (
	TemplateContext struct {
		scope *templateScope
	}
	templateScope struct {
		parent    *templateScope
		variables map[string]interface{}
	}
)

func NewRootTemplateContext() *TemplateContext {
	return &TemplateContext{&templateScope{nil, make(map[string]interface{})}}
}

func (c *TemplateContext) PushScope() {
	c.scope = c.scope.child()
}

func (c *TemplateContext) PopScope() {
	c.scope = c.scope.parent
}

func (c *TemplateContext) Get() map[string]interface{} {
	return c.scope.get()
}

func (c *TemplateContext) Set(key string, value interface{}) {
	c.scope.set(key, value)
}

func (c *templateScope) child() *templateScope {
	return &templateScope{c, make(map[string]interface{})}
}

func (c *templateScope) set(key string, value interface{}) {
	c.variables[key] = value
}

func (c *templateScope) get() map[string]interface{} {
	if c.parent != nil {
		vars := make(map[string]interface{})
		parentVars := c.parent.get()
		for k, v := range parentVars {
			vars[k] = v
		}
		for k, v := range c.variables {
			vars[k] = v
		}
		return vars
	}
	return c.variables
}
