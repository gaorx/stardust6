package sdcodegen

import (
	"fmt"
	"github.com/gaorx/stardust6/sdslog"
	"io/fs"
	"log/slog"
)

// Logger 生成代码时的可以打印的日志记录器
type Logger interface {
	// Log 记录一条日志
	Log(msg string)
	// LogWrite 记录一个文件写入操作
	LogWrite(name string, mode fs.FileMode, err error)
	// LogDiscard 记录一个文件丢弃操作
	LogDiscard(name string, mode fs.FileMode)
}

// NoLog 创建一个不记录日志的Logger
func NoLog() Logger {
	return nologLogger{}
}

// Stdout 创建一个输出到标准输出的Logger
func Stdout() Logger {
	return defaultLogger{}
}

// Slog 创建一个输出到slog的Logger
func Slog(logger *slog.Logger) Logger {
	if logger == nil {
		logger = slog.Default()
	}
	return slogLogger{logger: logger}
}

// SetLogger 创建一个中间件，用于设置Logger
func SetLogger(logger Logger) Middleware {
	return func(c *Context, next Handler) {
		c.SetLogger(logger)
		next(c)
	}
}

type defaultLogger struct{}

func (l defaultLogger) Log(msg string) {
	fmt.Println(msg)
}

func (l defaultLogger) LogWrite(name string, mode fs.FileMode, err error) {
	if err != nil {
		mode.Perm()
		fmt.Printf("ERROR  %s %v\n", name, err.Error())
		return
	}
	fmt.Printf("WRITE  %s\n", name)
}

func (l defaultLogger) LogDiscard(name string, mode fs.FileMode) {
	fmt.Printf("DISCARD %s\n", name)
}

type slogLogger struct {
	logger *slog.Logger
}

func (l slogLogger) Log(msg string) {
	l.logger.Info(msg)
}

func (l slogLogger) LogWrite(name string, mode fs.FileMode, err error) {
	if err != nil {
		l.logger.Error(fmt.Sprintf("ERROR  %s", name), sdslog.E(err))
		return
	}
	l.logger.Info(fmt.Sprintf("WRITE  %s", name))
}

func (l slogLogger) LogDiscard(name string, mode fs.FileMode) {
	l.logger.Info(fmt.Sprintf("DISCARD %s", name))
}

type nologLogger struct{}

func (l nologLogger) Log(msg string)                                    {}
func (l nologLogger) LogWrite(name string, mode fs.FileMode, err error) {}
func (l nologLogger) LogDiscard(name string, mode fs.FileMode)          {}
