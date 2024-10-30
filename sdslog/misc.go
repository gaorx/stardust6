package sdslog

import (
	"github.com/samber/lo"
	slogmulti "github.com/samber/slog-multi"
)

// Wrap 将多个middleware串联起来，施加到一个Handler上
func Wrap(h Handler, middlewares ...Middleware) Handler {
	middlewares = lo.Filter(middlewares, func(m Middleware, _ int) bool {
		return m != nil
	})
	if len(middlewares) <= 0 {
		return h
	}
	return slogmulti.Pipe(middlewares...).Handler(h)
}

// IsTTY 判断一个输出日志的文件名是否是在终端上输出，通常用来配合TextHandler的Pretty属性
func IsTTY(file string) bool {
	if isStdoutFile(file) {
		return isStdoutTTY()
	} else if isStderrFile(file) {
		return isStderrTTY()
	} else {
		return false
	}
}
