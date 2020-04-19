package domain

import (
	"fmt"
)

type (
	Result struct {
		JobID     string
		Label     string
		Link      string
		Previous  string
		Current   string
		ValueType ValueType
	}
)

func NewResult(info *JobInfo, previous, current string, valueType ValueType) *Result {
	return &Result{
		JobID:     info.ID,
		Label:     info.Label,
		Link:      info.Link,
		Previous:  previous,
		Current:   current,
		ValueType: valueType,
	}
}

func (r *Result) Diff() (DiffResult, error) {
	switch r.ValueType {
	case ValueTypeString:
		return DiffString(r.Previous, r.Current), nil
	case ValueTypeJSONObject:
		return DiffJSONObject(r.Previous, r.Current)
	case ValueTypeJSONArray:
		return DiffJSONArray(r.Previous, r.Current)
	default:
		return nil, fmt.Errorf("unsupported value type: %v", r.ValueType)
	}
}
