package logger

import (
	"log/slog"
	"os"
)

// Setup устанавливает настройки логирования
func Setup(env string) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: getLogLevel(env),
		//AddSource: env != "prod",
	}

	if env == "prod" {
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}

// getLogLevel возвращает уровень логирования для текущего окружения
func getLogLevel(env string) slog.Level {
	switch env {
	case "prod":
		return slog.LevelInfo
	default:
		return slog.LevelDebug
	}
}
