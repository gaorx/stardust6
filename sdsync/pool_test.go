package sdsync

import (
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	is := assert.New(t)

	pool, err := NewPool(10, nil)
	is.NoError(err)
	a := make([]int, 100)
	start := time.Now()
	err = ForEach(lo.Range(100), func(v, i int) {
		err0 := pool.Do(func() {
			a[i] = v * 2
			time.Sleep(200 * time.Millisecond)
		})
		is.NoError(err0)
	}, 10)
	elapsed := time.Since(start).Milliseconds()
	is.NoError(err)
	for i := 0; i < 100; i++ {
		is.Equal(i*2, a[i])
	}
	is.True(elapsed >= 1900 && elapsed <= 2100)
}
