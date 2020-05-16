package template

type (
	TemplateString string
)

func (t TemplateString) Evaluate(ctx *TemplateContext) (string, error) {
	tmpl, err := Parse(string(t))
	if err != nil {
		return "", err
	}
	return tmpl.Evaluate(ctx)
}
