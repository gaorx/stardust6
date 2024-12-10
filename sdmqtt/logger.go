package sdmqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
)

// Loggers MQTT日志配置
type Loggers struct {
	Debug    mqtt.Logger
	Error    mqtt.Logger
	Warn     mqtt.Logger
	Critical mqtt.Logger
}

// SetLoggers 设置MQTT全局日志
func SetLoggers(loggers Loggers) {
	if loggers.Debug != nil {
		mqtt.DEBUG = loggers.Debug
	}
	if loggers.Error != nil {
		mqtt.ERROR = loggers.Error
	}
	if loggers.Warn != nil {
		mqtt.WARN = loggers.Warn
	}
	if loggers.Critical != nil {
		mqtt.CRITICAL = loggers.Critical
	}
}

// SetSlog 设置MQTT全局日志为slog
func SetSlog(logger *slog.Logger) {
	SetLoggers(Loggers{
		Debug:    SlogOf(logger, slog.LevelDebug),
		Error:    SlogOf(logger, slog.LevelError),
		Warn:     SlogOf(logger, slog.LevelWarn),
		Critical: SlogOf(logger, slog.LevelWarn),
	})
}
