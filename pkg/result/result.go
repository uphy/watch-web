package result

import (
	"fmt"

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

func (r *Result) Diff() (*DiffResult, error) {
	switch r.ValueType {
	case value.ValueTypeString:
		return Diff(r.Previous, r.Current), nil
	case value.ValueTypeJSONObject:
		return DiffJSONObject(r.Previous, r.Current)
	case value.ValueTypeJSONArray:
		return DiffJSONArray(r.Previous, r.Current)
	default:
		return nil, fmt.Errorf("unsupported value type: %v", r.ValueType)
	}
}
