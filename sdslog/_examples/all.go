package main

import (
	"context"
	"fmt"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdslog"
	"io/fs"
	"log/slog"
	"strings"
	"time"
)

func main() {
	exampleNew()
}
func separator() {
	fmt.Println(strings.Repeat("-", 60))
}

func exampleNew() {
	fmt.Println("EXAMPLE: NEW")
	separator()

	// 构建一个只有一个text handler的logger，日志等级为Debug，输出到标准输出，格式为便于阅读的格式
	logger1 := sdslog.New([]sdslog.Handler{
		sdslog.TextFile(slog.LevelDebug, "stdout", true),
	}, nil)
	logger1.Info("message", "k1", "v1")
	separator()

	// 构建一个只有一个text handler的logger，日志等级为Debug，输出到标准输出
	// 只有输出到终端上，则为便于阅读的格式，否则为plain文本
	logger2 := sdslog.New([]sdslog.Handler{
		sdslog.TextFile(slog.LevelDebug, "stdout", sdslog.IsTTY("stdout")),
	}, nil)
	logger2.Info("message", "k1", "v1")
	separator()

	// 构建一个只有一个json handler的logger，日志等级为Debug，输出到指定log文件
	logger3 := sdslog.New([]sdslog.Handler{
		sdslog.JsonFile(slog.LevelDebug, "stdout"),
	}, nil)
	logger3.Info("message", "k1", "v1")
	separator()

	// 构建一个有双路输出的logger，可以同时输出到标准输出和指定log文件
	logger4 := sdslog.New([]sdslog.Handler{
		sdslog.TextFile(slog.LevelDebug, "stdout", true),
		sdslog.JsonFile(slog.LevelDebug, "stdout"),
	}, nil)
	logger4.Info("message", "k1", "v1")
	separator()

	// 构建一个有自定义输出函数的logger
	logger5 := sdslog.New([]sdslog.Handler{
		sdslog.HandleRecord(func(ctx context.Context, record sdslog.Record) error {
			fmt.Println("******", record.Message)
			return nil
		}),
	}, nil)
	logger5.Info("message", "k1", "v1")
	separator()

	// 构建一个logger，它可以打印错误的stacktrace或者展开错误的attr
	logger6 := sdslog.New([]sdslog.Handler{
		sdslog.TextFile(slog.LevelDebug, "stdout", true),
	}, []sdslog.Middleware{
		sdslog.ExpandError(&sdslog.ExpandErrorOptions{Stack: true}),
	})
	logger6.Error("error", sdslog.E(newError2()))
	separator()

	// 构建一个logger，它可以
	logger7 := sdslog.New([]sdslog.Handler{
		sdslog.TextFile(slog.LevelDebug, "stdout", true),
	}, []sdslog.Middleware{
		sdslog.Format(
			// 格式化k1的值，将某个隐私的值替换为"******"
			sdslog.FormatByKey("privacy_key", func(value sdslog.Value) sdslog.Value {
				return slog.StringValue("******")
			}),
			sdslog.FormatTime(time.Kitchen, nil), // 识别Attr中的time.Time类型，格式化为指定的时间格式
		),
	})
	logger7.Info("message", "privacy_key", "<phone_number>")
	separator()

	// 构建一个logger，它可以拦截所有的日志处理函数，在打印每行日志之前为每行日志的消息加入前缀"****"
	logger8 := sdslog.New([]sdslog.Handler{
		sdslog.TextFile(slog.LevelDebug, "stdout", true),
	}, []sdslog.Middleware{
		sdslog.InterceptHandle(func(ctx context.Context, record sdslog.Record, next func(context.Context, sdslog.Record) error) error {
			record.Message = "****" + record.Message
			return next(ctx, record)
		}),
	})
	logger8.Info("message", "k1", "v1")
	separator()
}

func newError1() error {
	return sderr.With("k1", "v1").Wrapf(fs.ErrExist, "error1")
}

func newError2() error {
	return sderr.With("k2", "v2").Wrapf(newError1(), "error2")
}
