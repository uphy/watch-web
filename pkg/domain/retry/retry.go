package retry

import (
	"math/rand"
	"time"
)

type (
	RetryableFunc func(ctx *RetryContext) error
	Retrier       struct {
		initialInterval     float64
		multiplier          float64
		randomizationFactor float64
		maxInterval         float64
		retry               int
	}
	RetryContext struct {
		Retried int
	}
	Builder struct {
		initialInterval     float64
		multiplier          float64
		randomizationFactor float64
		maxInterval         float64
		retry               int
	}
)

func NewBuilder(retry int) *Builder {
	return &Builder{
		initialInterval:     1.,
		multiplier:          2,
		randomizationFactor: 0.5,
		maxInterval:         0,
		retry:               retry,
	}
}

func (b *Builder) InitialInterval(initialInterval float64) *Builder {
	b.initialInterval = initialInterval
	return b
}

func (b *Builder) Multiplier(multiplier float64) *Builder {
	b.multiplier = multiplier
	return b
}

func (b *Builder) RandomizationFactor(randomizationFactor float64) *Builder {
	b.randomizationFactor = randomizationFactor
	return b
}

func (b *Builder) MaxInterval(maxInterval float64) *Builder {
	b.maxInterval = maxInterval
	return b
}

func (b *Builder) Build() *Retrier {
	return &Retrier{
		retry:               b.retry,
		initialInterval:     b.initialInterval,
		multiplier:          b.multiplier,
		randomizationFactor: b.randomizationFactor,
		maxInterval:         b.maxInterval,
	}
}

func (r *Retrier) Run(task RetryableFunc) error {
	var err error
	ctx := new(RetryContext)
	exponentialBackoff := r.initialInterval
	for {
		err = task(ctx)
		if err == nil {
			break
		}

		maxJitter := exponentialBackoff * r.randomizationFactor
		jitter := rand.Float64() * maxJitter
		wait := float64(time.Second) * (exponentialBackoff + jitter)
		if r.maxInterval > 0 && wait > r.maxInterval {
			wait = r.maxInterval
		} else {
			exponentialBackoff *= r.multiplier
		}
		time.Sleep(time.Duration(wait))

		ctx.Retried++

		if ctx.Retried >= r.retry {
			break
		}
	}
	return err
}
