package sdsemver

import (
	"github.com/gaorx/stardust6/sderr"
)

var (
	// ErrParse 解析字符串错误
	ErrParse = sderr.Sentinel("parse semver failed")
)
