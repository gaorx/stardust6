package sdbun

import (
	"github.com/gaorx/stardust6/sdsql"
)

// IsNotFoundErr 判断是否为NotFound语义的错误
func IsNotFoundErr(err error) bool {
	return sdsql.IsNotFoundErr(err)
}

// IsReadonlyErr 判断是否为只读错误
func IsReadonlyErr(err error) bool {
	return sdsql.IsReadonlyErr(err)
}
