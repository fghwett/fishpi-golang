package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
)

func init() {
	New()
}

func New() *slog.Logger {

	// 创建 lumberjack.Logger 对象，配置日志文件路径和其他属性
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "_tmp/log.log",
		MaxSize:    5,    // 每个日志文件的最大尺寸，单位为 MB
		MaxBackups: 3,    // 保留的旧日志文件的最大数量
		MaxAge:     30,   // 保留的旧日志文件的最大天数
		Compress:   true, // 是否压缩旧日志文件
	}

	logger := slog.New(slog.NewTextHandler(io.MultiWriter(lumberjackLogger), &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	return logger
}

func RecordMessage(message interface{}) {
	slog.Debug("message", slog.Any("message", message))
}
