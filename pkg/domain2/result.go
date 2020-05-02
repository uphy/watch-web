package domain2

type (
	Result struct {
		JobID    string
		Label    string
		Link     string
		Previous ItemList
		Current  ItemList
	}
)

func NewResult(jobID string, label string, link string, prev, current ItemList) *Result {
	return &Result{
		JobID:    jobID,
		Label:    label,
		Link:     link,
		Previous: prev,
		Current:  current,
	}
}

func (r *Result) Diff() *Updates {
	return CompareItemList(r.Previous, r.Current)
}
