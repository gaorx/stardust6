package sdsync

import (
	"github.com/samber/lo"
	"sync"
)

// Do 并发执行一组操作，numConcurrency是并发数量，如果<=0则表示所有actions都并发执行
func Do(actions []func(), numConcurrency int) error {
	numActions := len(actions)
	if numActions == 0 {
		return nil
	}
	if numConcurrency <= 0 {
		var wg sync.WaitGroup
		for _, f := range actions {
			wg.Add(1)
			go func(f func()) {
				defer wg.Done()
				_ = lo.Try0(f)
			}(f)
		}
		wg.Wait()
		return nil
	} else {
		if numConcurrency > numActions {
			numConcurrency = numActions
		}
		pool, err := NewPool(numConcurrency, &PoolOptions{
			PreAlloc: true,
		})
		if err != nil {
			return err
		}
		defer func() { _ = pool.Close() }()
		var wg sync.WaitGroup
		for _, f := range actions {
			f1 := f
			wg.Add(1)
			err := pool.Submit(func() {
				defer wg.Done()
				_ = lo.Try0(f1)
			})
			if err != nil {
				return err
			}
		}
		wg.Wait()
		return nil
	}
}

// ForEach 并发对一个slice的每个元素进行并发操作，numConcurrency是并发数量，如果<=0则表示所有元素都并发执行
func ForEach[T any](l []T, action func(int, T), numConcurrency int) error {
	if action == nil {
		action = func(int, T) {}
	}
	actions := make([]func(), 0, len(l))
	for i, v := range l {
		i0, v0 := i, v
		actions = append(actions, func() {
			action(i0, v0)
		})
	}
	return Do(actions, numConcurrency)
}
