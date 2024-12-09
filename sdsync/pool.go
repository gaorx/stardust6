package sdsync

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/panjf2000/ants/v2"
	"github.com/samber/lo"
	"sync"
)

// Pool 协程池
type Pool struct {
	pool *ants.Pool
}

// PoolOptions 池选项
type PoolOptions = ants.Options

var (
	ErrInvalidPoolExpiry   = ants.ErrInvalidPoolExpiry
	ErrLackPoolFunc        = ants.ErrLackPoolFunc
	ErrPoolClosed          = ants.ErrPoolClosed
	ErrPoolOverload        = ants.ErrPoolOverload
	ErrInvalidPreAllocSize = ants.ErrInvalidPreAllocSize
)

// NewPool 创建一个协程池
func NewPool(size int, opts *PoolOptions) (*Pool, error) {
	var antsOpts []ants.Option
	if opts != nil {
		antsOpts = append(antsOpts, ants.WithOptions(*opts))
	}
	p, err := ants.NewPool(size, antsOpts...)
	if err != nil {
		return nil, sderr.Wrapf(err, "create ants pool error")
	}
	return &Pool{pool: p}, nil
}

// NumFree 获取协程池中还空闲的协程数量
func (p *Pool) NumFree() int {
	return p.pool.Free()
}

// NumCap 获取协程池中协程容量
func (p *Pool) NumCap() int {
	return p.pool.Cap()
}

// NumRunning 获取协程池中正在运行的协程数量
func (p *Pool) NumRunning() int {
	return p.pool.Running()
}

// Close 关闭协程池
func (p *Pool) Close() error {
	p.pool.Release()
	return nil
}

// Submit 提交一个操作到协程池，不等待操作完成立刻返回
func (p *Pool) Submit(action func()) error {
	if action == nil {
		return nil
	}
	err := p.pool.Submit(action)
	return sderr.Wrapf(err, "submit action error")
}

// Do 向协程池提交一个操作，等待操作完成后返回
func (p *Pool) Do(action func()) error {
	if action == nil {
		return nil
	}
	var wg sync.WaitGroup
	wg.Add(1)
	err := p.pool.Submit(func() {
		defer wg.Done()
		_ = lo.Try0(action)
	})
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

// Wrap 包装一个操作，返回一个新的操作，新的操作会在协程池中执行
func (p *Pool) Wrap(action func()) func() {
	if action == nil {
		return nil
	}
	return func() {
		_ = p.Do(action)
	}
}
