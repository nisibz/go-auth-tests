package logger

import (
	"log/slog"
	"os"

	"github.com/nisibz/go-auth-tests/internal/adapter/config"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

// logger is the default logger used by the application
var logger *slog.Logger

// Set sets the logger configuration based on the environment
func Set(config *config.App) {
	logger = slog.New(
		slog.NewTextHandler(os.Stderr, nil),
	)

	logRotate := &lumberjack.Logger{
		Filename:   "log/app.log",
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	logger = slog.New(
		slogmulti.Fanout(
			slog.NewJSONHandler(logRotate, nil),
			slog.NewTextHandler(os.Stderr, nil),
		),
	)

	slog.SetDefault(logger)
}
