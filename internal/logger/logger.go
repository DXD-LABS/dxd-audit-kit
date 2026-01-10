package logger

import (
	"log/slog"
	"os"
)

var defaultLogger *slog.Logger

func init() {
	// Mặc định sử dụng JSON handler cho production readiness
	defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

func Error(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, slog.Any("error", err))
	}
	defaultLogger.Error(msg, args...)
}

func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}
