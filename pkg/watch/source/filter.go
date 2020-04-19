package source

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/uphy/watch-web/pkg/domain"
)

type (
	FilterSource struct {
		source  domain.Source
		filters []domain.Filter
	}
)

func NewFilterSource(source domain.Source, filters []domain.Filter) *FilterSource {
	return &FilterSource{source, filters}
}

func (f *FilterSource) Fetch(ctx *domain.JobContext) (domain.Value, error) {
	v, err := f.source.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	ctx.Log.WithField("source", v).Debug("Start filter chain.")
	for _, filter := range f.filters {
		filtered, err := filter.Filter(ctx, v)
		if err != nil {
			ctx.Log.WithFields(logrus.Fields{
				"filter": fmt.Sprintf("%#v", filter),
			}).Debug("Failed to filtered value.")
			return nil, err
		}
		ctx.Log.WithFields(logrus.Fields{
			"filter":   fmt.Sprintf("%#v", filter),
			"filtered": filtered,
		}).Debug("Filtered value.")
		v = filtered
	}
	return v, nil
}

func (f *FilterSource) String() string {
	return fmt.Sprintf("Filter[source=%v, filters=%v]", f.source, f.filters)
}
