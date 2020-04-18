package result

import (
	"github.com/uphy/watch-web/pkg/value"
)

type (
	Result struct {
		JobID     string
		Label     string
		Link      string
		Previous  string
		Current   string
		ValueType value.ValueType
	}
)

func New(jobId string, label string, link string, previous string, current string, valueType value.ValueType) *Result {
	return &Result{
		JobID:     jobId,
		Label:     label,
		Link:      link,
		Previous:  previous,
		Current:   current,
		ValueType: valueType,
	}
}

func (r *Result) Diff() *Diff {
	return diff(r.Previous, r.Current)
}
