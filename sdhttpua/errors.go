package sdhttpua

import (
	"github.com/gaorx/stardust6/sderr"
)

var (
	// ErrParse 解析User-Agent错误
	ErrParse = sderr.Sentinel("parse User-Agent error")
)
