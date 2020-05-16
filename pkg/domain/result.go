package domain

import "github.com/uphy/watch-web/pkg/domain/value"

type (
	Result struct {
		JobID    string
		Label    string
		Link     string
		Previous value.ItemList
		Current  value.ItemList
	}
)

func NewResult(info *JobInfo, previous, current value.ItemList) *Result {
	return &Result{
		JobID:    info.ID,
		Label:    info.Label,
		Link:     info.Link,
		Previous: previous,
		Current:  current,
	}
}

func (r *Result) Diff() value.Updates {
	return value.CompareItemList(r.Previous, r.Current)
}
