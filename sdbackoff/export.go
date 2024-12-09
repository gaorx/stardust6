package sdbackoff

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
	"time"
)

type (
	// BackOff 补偿策略
	BackOff = backoff.BackOff
	// Ticker Ticker
	Ticker = backoff.Ticker
)

// ExponentialOptions 指数选项
type ExponentialOptions struct {
	InitialInterval     time.Duration
	RandomizationFactor float64
	Multiplier          float64
	MaxInterval         time.Duration
	MaxElapsedTime      time.Duration
}

// Exponential 指数补偿
func Exponential(opts *ExponentialOptions) BackOff {
	opts1 := lo.FromPtr(opts)
	if opts1.InitialInterval <= 0 {
		opts1.InitialInterval = backoff.DefaultInitialInterval
	}
	if opts1.RandomizationFactor <= 0.0 {
		opts1.RandomizationFactor = backoff.DefaultRandomizationFactor
	}
	if opts1.Multiplier <= 0.0 {
		opts1.Multiplier = backoff.DefaultMultiplier
	}
	if opts1.MaxInterval <= 0 {
		opts1.MaxInterval = backoff.DefaultMaxInterval
	}
	if opts1.MaxElapsedTime <= 0 {
		opts1.MaxElapsedTime = backoff.DefaultMaxElapsedTime
	}
	b := &backoff.ExponentialBackOff{
		InitialInterval:     opts1.InitialInterval,
		RandomizationFactor: opts1.RandomizationFactor,
		Multiplier:          opts1.Multiplier,
		MaxInterval:         opts1.MaxInterval,
		MaxElapsedTime:      opts1.MaxElapsedTime,
		Clock:               backoff.SystemClock,
	}
	b.Reset()
	return b
}

// Stop stop补偿
func Stop() BackOff {
	return &backoff.StopBackOff{}
}

// Zero 零补偿
func Zero() BackOff {
	return &backoff.ZeroBackOff{}
}

// Const 常量补偿
func Const(d time.Duration) BackOff {
	if d > 0 {
		return backoff.NewConstantBackOff(d)
	} else {
		return &backoff.ZeroBackOff{}
	}
}

// TickerOf 通过补偿创建一个ticker
func TickerOf(b BackOff) *Ticker {
	return backoff.NewTicker(b)
}

// Retry 通过补偿进行重试
func Retry(b BackOff, action func() error) error {
	err := backoff.Retry(action, b)
	return sderr.Wrapf(err, "sdbackoff retry error")
}
