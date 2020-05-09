package domain

type (
	Result struct {
		JobID    string
		Label    string
		Link     string
		Previous ItemList
		Current  ItemList
	}
)

func NewResult(info *JobInfo, previous, current ItemList) *Result {
	return &Result{
		JobID:    info.ID,
		Label:    info.Label,
		Link:     info.Link,
		Previous: previous,
		Current:  current,
	}
}

func (r *Result) Diff() Updates {
	return CompareItemList(r.Previous, r.Current)
}
