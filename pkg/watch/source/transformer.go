package source

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/uphy/watch-web/pkg/domain"
)

type (
	TransformerSource struct {
		source       domain.Source
		transformers []domain.Transformer
	}
)

func NewTransformerSource(source domain.Source, transformers []domain.Transformer) *TransformerSource {
	return &TransformerSource{source, transformers}
}

func (f *TransformerSource) Fetch(ctx *domain.JobContext) (domain.Value, error) {
	v, err := f.source.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	ctx.Log.WithField("source", v).Debug("Start transformer chain.")
	for _, transformer := range f.transformers {
		transformed, err := transformer.Transform(ctx, v)
		if err != nil {
			ctx.Log.WithFields(logrus.Fields{
				"transformer": fmt.Sprintf("%#v", transformer),
			}).Debug("Failed to transform value.")
			return nil, err
		}
		ctx.Log.WithFields(logrus.Fields{
			"transformer": fmt.Sprintf("%#v", transformer),
			"transformed": transformed,
		}).Debug("Transformed value.")
		v = transformed
	}
	return v, nil
}

func (f *TransformerSource) String() string {
	return fmt.Sprintf("Transformer[source=%v, transformers=%v]", f.source, f.transformers)
}
