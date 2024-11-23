package sderr

import (
	stderr "errors"
)

// Sentinel 创建一个最底层的无须栈信息的错误，用于定义全局错误
func Sentinel(text string) error {
	return stderr.New(text)
}
