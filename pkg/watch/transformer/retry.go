package transformer

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/domain/retry"
	"github.com/uphy/watch-web/pkg/domain/value"
)

type (
	TransformerWithRetry struct {
		transformer domain.Transformer
		retrier     *retry.Retrier
	}
)

func NewRetryTransform(transformer domain.Transformer, retrier *retry.Retrier) *TransformerWithRetry {
	return &TransformerWithRetry{transformer, retrier}
}

func (s *TransformerWithRetry) Transform(ctx *domain.JobContext, v value.Value) (value.Value, error) {
	if s.retrier == nil {
		return s.transformer.Transform(ctx, v)
	}
	var result value.Value
	err := s.retrier.Run(func(retryContext *retry.RetryContext) error {
		var err error
		result, err = s.transformer.Transform(ctx, v)
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *TransformerWithRetry) String() string {
	var retry = ""
	if s.retrier != nil {
		retry = fmt.Sprint(*s.retrier)
	}
	return fmt.Sprintf("Transform[source=%v, retry=%v]", s.transformer, retry)
}
