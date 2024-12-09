package sdsync

import (
	"github.com/gaorx/stardust6/sdtime"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestForEach(t *testing.T) {
	is := assert.New(t)

	const (
		threads = 100
		n       = 200
		sleepMS = 10
	)
	var mtx sync.Mutex
	counter := 0
	start := sdtime.NowUnixMS()
	err := ForEach(lo.Range(threads), func(_, _ int) {
		for i := 0; i < n; i++ {
			sdtime.SleepMS(sleepMS)
			Lock(&mtx, func() {
				counter += 1
			})
		}
	}, 0)
	elapsed := sdtime.NowUnixMS() - start
	is.NoError(err)
	is.True(elapsed >= sleepMS*n && elapsed <= 2*sleepMS*n)
	is.Equal(threads*n, counter)
}
