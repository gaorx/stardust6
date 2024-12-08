package sdcache

import (
	"context"
	"github.com/gaorx/stardust6/sderr"
	"time"
)

// Double 双层缓存，先从L1读取，如果不存在则从L2读取并写入L1
func Double[K Key, V any](l1, l2 Cache[K, V]) Cache[K, V] {
	return doubleLayer[K, V]{l1: l1, l2: l2}
}

type doubleLayer[K Key, V any] struct {
	l1, l2 Cache[K, V]
}

var _ Cache[string, string] = doubleLayer[string, string]{}

func (d doubleLayer[K, V]) Clear(ctx context.Context) error {
	l1, l2 := d.l1, d.l2
	if l1 != nil && l2 != nil {
		err2 := l2.Clear(ctx)
		err1 := l1.Clear(ctx)
		return joinErr(err2, err1)
	} else if l1 != nil {
		return l1.Clear(ctx)
	} else if l2 != nil {
		return l2.Clear(ctx)
	} else {
		return sderr.Newf("double cache is empty")
	}
}

func (d doubleLayer[K, V]) Get(ctx context.Context, key K) (V, error) {
	l1, l2 := d.l1, d.l2
	if l1 != nil && l2 != nil {
		return l1.GetOrLoad(ctx, key, func(ctx context.Context, k K) (V, bool, error) {
			v0, err0 := l2.Get(ctx, k)
			if err0 != nil {
				var zero V
				if sderr.Is(err0, ErrNotFound) {
					return zero, false, nil
				}
				return zero, false, err0
			}
			return v0, true, nil
		}, nil)
	} else if l1 != nil {
		return l1.Get(ctx, key)
	} else if l2 != nil {
		return l2.Get(ctx, key)
	} else {
		var zero V
		return zero, sderr.Newf("double cache is empty")
	}
}

func (d doubleLayer[K, V]) GetTTL(ctx context.Context, key K) (time.Duration, error) {
	l1, l2 := d.l1, d.l2
	if l1 != nil {
		return l1.GetTTL(ctx, key)
	} else if l2 != nil {
		return l2.GetTTL(ctx, key)
	} else {
		return 0, sderr.Newf("double cache is empty")
	}
}

func (d doubleLayer[K, V]) Put(ctx context.Context, key K, val V, opts *PutOptions) error {
	l1, l2 := d.l1, d.l2
	if l1 != nil && l2 != nil {
		err2 := l2.Put(ctx, key, val, opts)
		err1 := l1.Put(ctx, key, val, nil)
		return joinErr(err1, err2)
	} else if l1 != nil {
		return l1.Put(ctx, key, val, opts)
	} else if l2 != nil {
		return l2.Put(ctx, key, val, opts)
	} else {
		return sderr.Newf("double cache is empty")
	}
}

func (d doubleLayer[K, V]) Delete(ctx context.Context, key K) error {
	l1, l2 := d.l1, d.l2
	if l1 != nil && l2 != nil {
		err2 := l2.Delete(ctx, key)
		err1 := l1.Delete(ctx, key)
		return joinErr(err1, err2)
	} else if l1 != nil {
		return l1.Delete(ctx, key)
	} else if l2 != nil {
		return l2.Delete(ctx, key)
	} else {
		return sderr.Newf("double cache is empty")
	}
}

func (d doubleLayer[K, V]) GetOrLoad(ctx context.Context, key K, loader Loader[K, V], opts *PutOptions) (V, error) {
	l1, l2 := d.l1, d.l2
	if l1 != nil && l2 != nil {
		return l1.GetOrLoad(ctx, key, func(ctx context.Context, k K) (V, bool, error) {
			v0, err0 := l2.GetOrLoad(ctx, k, loader, opts)
			if err0 != nil {
				var zero V
				if sderr.Is(err0, ErrNotFound) {
					return zero, false, nil
				}
				return zero, false, err0
			}
			return v0, true, nil
		}, nil)
	} else if l1 != nil {
		return l1.GetOrLoad(ctx, key, loader, opts)
	} else if l2 != nil {
		return l2.GetOrLoad(ctx, key, loader, opts)
	} else {
		var zero V
		return zero, sderr.Newf("double cache is empty")
	}
}

func joinErr(err1, err2 error) error {
	if err1 == nil && err2 == nil {
		return nil
	} else if err1 != nil && err2 != nil {
		return sderr.Join(err1, err2)
	} else if err1 != nil {
		return err1
	} else {
		return err2
	}
}
