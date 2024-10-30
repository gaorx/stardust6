package sdslog

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func Example_new() {
	// 构建一个只有一个text handler的logger，日志等级为Debug，输出到标准输出，格式为便于阅读的格式
	logger1 := New([]Handler{
		TextFile(slog.LevelDebug, "stdout", true),
	}, nil)
	logger1.Info("message", "k1", "v1")

	// 构建一个只有一个text handler的logger，日志等级为Debug，输出到标准输出
	// 只有输出到终端上，则为便于阅读的格式，否则为plain文本
	logger2 := New([]Handler{
		TextFile(slog.LevelDebug, "stdout", IsTTY("stdout")),
	}, nil)
	logger2.Info("message", "k1", "v1")

	// 构建一个只有一个json handler的logger，日志等级为Debug，输出到指定log文件
	logger3 := New([]Handler{
		JsonFile(slog.LevelDebug, "/path/to/logfile"),
	}, nil)
	logger3.Info("message", "k1", "v1")

	// 构建一个有双路输出的logger，可以同时输出到标准输出和指定log文件
	logger4 := New([]Handler{
		TextFile(slog.LevelDebug, "stdout", true),
		JsonFile(slog.LevelDebug, "/path/to/logfile"),
	}, nil)
	logger4.Info("message", "k1", "v1")

	// 构建一个有自定义输出函数的logger
	logger5 := New([]Handler{
		HandleRecord(func(ctx context.Context, record Record) error {
			fmt.Println("******", record.Message)
			return nil
		}),
	}, nil)
	logger5.Info("message", "k1", "v1")

	// 构建一个logger，它可以将识别sderr.Error，将起中的Attr展开到Record的group中
	logger6 := New([]Handler{
		TextFile(slog.LevelDebug, "stdout", true),
	}, []Middleware{
		ExpandSderrError,
	})
	logger6.Info("message", "k1", "v1")

	// 构建一个logger，它可以
	logger7 := New([]Handler{
		TextFile(slog.LevelDebug, "stdout", true),
	}, []Middleware{
		Format(
			// 格式化k1的值，将某个隐私的值替换为"******"
			FormatByKey("privacy_key", func(value Value) Value {
				return slog.StringValue("******")
			}),
			FormatTime(time.Kitchen, nil), // 识别Attr中的time.Time类型，格式化为指定的时间格式
		),
	})
	logger7.Info("message", "privacy_key", "<phone_number>")

	// 构建一个logger，它可以拦截所有的日志处理函数，在打印每行日志之前为每行日志的消息加入前缀"****"
	logger8 := New([]Handler{
		TextFile(slog.LevelDebug, "stdout", true),
	}, []Middleware{
		InterceptHandle(func(ctx context.Context, record Record, next func(context.Context, Record) error) error {
			record.Message = "****" + record.Message
			return next(ctx, record)
		}),
	})
	logger8.Info("message", "k1", "v1")
}
