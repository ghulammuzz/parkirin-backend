package log

import (
	"log/slog"

	"github.com/grafana/loki-client-go/loki"
	slogloki "github.com/samber/slog-loki/v3"
)

var Logger *slog.Logger

func InitLogger(profile string, lokiClient *loki.Client) {
	SetProfileLog(profile, lokiClient)
}

func SetProfileLog(profile string, lokiClient *loki.Client) {
	var level slog.Leveler

	switch profile {
	case "dev":
		level = slog.LevelDebug
	case "prod":
		level = slog.LevelInfo
	default:
		level = slog.LevelInfo
	}

	// opts := &slog.HandlerOptions{
	// 	Level: level,
	// }

	// handler := slog.NewTextHandler(os.Stdout, opts)

	Logger = slog.New(slogloki.Option{Level: level, Client: lokiClient}.NewLokiHandler())
	Logger = Logger.
		With("environment", profile).
		With("apps_name", "be-parkirin").
		With("release", "v1.1")
}

func Debug(msg string, args ...interface{}) {
	if Logger != nil {
		Logger.Debug(msg, args...)
	}
}

func Info(msg string, args ...interface{}) {
	if Logger != nil {
		Logger.Info(msg, args...)
	}
}

func Warn(msg string, args ...interface{}) {
	if Logger != nil {
		Logger.Warn(msg, args...)
	}
}

func Error(msg string, args ...interface{}) {
	if Logger != nil {
		Logger.Error(msg, args...)
	}
}
