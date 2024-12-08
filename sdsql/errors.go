package sdsql

import (
	"database/sql"
	"github.com/gaorx/stardust6/sderr"
)

var (
	// ErrReadonly 如果一个repository是只读的，调用了写操作，会返回这个错误
	ErrReadonly = sderr.Sentinel("readonly repository")
)

// IsNotFoundErr 判断是否是未找到错误
func IsNotFoundErr(err error) bool {
	return sderr.Is(err, sql.ErrNoRows)
}

// IsReadonlyErr 判断是否是只读错误
func IsReadonlyErr(err error) bool {
	return sderr.Is(err, ErrReadonly)
}
