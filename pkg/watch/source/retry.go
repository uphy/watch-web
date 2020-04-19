package source

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/uphy/watch-web/pkg/domain"
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
		retry       *int
	}
)

func NewRetrySource(src domain.Source, emptyAction *EmptyAction, retry *int) *SourceWithRetry {
	return &SourceWithRetry{src, emptyAction, retry}
}

func (s *SourceWithRetry) Fetch(ctx *domain.JobContext) (domain.Value, error) {
	if s.retry == nil {
		return s.fetch(ctx)
	}
	retry := *s.retry
	var err error
	for i := 0; i <= retry; i++ {
		var v domain.Value
		v, err = s.fetch(ctx)
		if err == nil {
			return v, nil
		}

		if i < retry {
			waitSeconds := 1 + int(math.Pow(2, float64(i)))
			time.Sleep(time.Second * time.Duration(waitSeconds))
		}
	}
	return nil, fmt.Errorf("too many retries: lastError=%w", err)
}

func (s *SourceWithRetry) fetch(ctx *domain.JobContext) (domain.Value, error) {
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
	if s.retry != nil {
		retry = fmt.Sprint(*s.retry)
	}
	return fmt.Sprintf("Retry[source=%v, retry=%v, emptyAction=%v]", s.source, retry, emptyAction)
}
