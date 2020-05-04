package domain

import (
	"github.com/uphy/watch-web/pkg/domain2"
)

type (
	Result struct {
		JobID    string
		Label    string
		Link     string
		Previous domain2.ItemList
		Current  domain2.ItemList
	}
)

func NewResult(info *JobInfo, previous, current domain2.ItemList) *Result {
	return &Result{
		JobID:    info.ID,
		Label:    info.Label,
		Link:     info.Link,
		Previous: previous,
		Current:  current,
	}
}

func (r *Result) Diff() domain2.Updates {
	return domain2.CompareItemList(r.Previous, r.Current)
}
