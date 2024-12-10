package sdmqtt

import (
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gaorx/stardust6/sdslog"
	"log/slog"
)

// SlogLogger slog日志适配
type SlogLogger struct {
	logger *slog.Logger
	level  slog.Level
}

var _ mqtt.Logger = SlogLogger{}

// SlogOf 创建slog日志适配
func SlogOf(logger *slog.Logger, level slog.Level) SlogLogger {
	if logger == nil {
		logger = sdslog.DiscardLogger
	}
	return SlogLogger{logger: logger, level: level}
}

// Println 实现 mqtt.Logger
func (s SlogLogger) Println(v ...any) {
	s.logger.Log(context.Background(), s.level, fmt.Sprint(v...))
}

// Printf 实现 mqtt.Logger
func (s SlogLogger) Printf(format string, v ...any) {
	s.logger.Log(context.Background(), s.level, fmt.Sprintf(format, v...))
}
