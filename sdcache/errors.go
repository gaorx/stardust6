package sdcache

import (
	"github.com/gaorx/stardust6/sderr"
)

var (
	// ErrNotFound 缓存key未找到返回的错误
	ErrNotFound = sderr.Sentinel("cache key not found")
)
