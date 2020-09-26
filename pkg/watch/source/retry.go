package source

import (
	"errors"
	"fmt"

	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/domain/retry"
	"github.com/uphy/watch-web/pkg/domain/value"
)

const (
	EmptyActionAccept EmptyAction = "accept"
	EmptyActionFail   EmptyAction = "fail"
)

type (
	EmptyAction     string
	SourceWithRetry struct {
		source      domain.Source
		emptyAction *EmptyAction
		retrier     *retry.Retrier
	}
)

func NewRetrySource(src domain.Source, emptyAction *EmptyAction, retrier *retry.Retrier) *SourceWithRetry {
	return &SourceWithRetry{src, emptyAction, retrier}
}

func (s *SourceWithRetry) Fetch(ctx *domain.JobContext) (value.Value, error) {
	if s.retrier == nil {
		return s.fetch(ctx)
	}
	var v value.Value
	err := s.retrier.Run(func(retryContext *retry.RetryContext) error {
		var err error
		v, err = s.fetch(ctx)
		return err
	})
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (s *SourceWithRetry) fetch(ctx *domain.JobContext) (value.Value, error) {
	v, err := s.source.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	if s.emptyAction == nil || *s.emptyAction == EmptyActionAccept {
		return v, nil
	}
	if *s.emptyAction == EmptyActionFail {
		if v.Empty() {
			return nil, errors.New("empty value")
		}
		return v, nil
	}
	return nil, fmt.Errorf("unsupported empty action: %s", *s.emptyAction)
}

func (s *SourceWithRetry) String() string {
	var emptyAction = ""
	var retry = ""
	if s.emptyAction != nil {
		emptyAction = string(*s.emptyAction)
	}
	if s.retrier != nil {
		retry = fmt.Sprint(*s.retrier)
	}
	return fmt.Sprintf("Retry[source=%v, retry=%v, emptyAction=%v]", s.source, retry, emptyAction)
}
