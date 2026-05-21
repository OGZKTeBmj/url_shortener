package logger

import (
	"io"
	"log/slog"
	"strings"
)

type Logger interface {
	Info(msg string, fields ...any)
	Debug(msg string, fields ...any)
	Error(msg string, fields ...any)

	With(fields ...any) Logger
}

type Config struct {
	AppName    string `env:"APP_NAME,required"`
	AppVersion string `env:"APP_VERSION,required"`
	Level      string `env:"LOG_LEVEL" envDefault:"info"`
	Env        string `env:"APP_ENV" envDefault:"prod"`
}

type SlogLogger struct {
	log *slog.Logger
}

func NewSlogLogger(w io.Writer, cfg Config) *SlogLogger {
	opts := &slog.HandlerOptions{
		Level: cfg.logLevel(),

		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return slog.Attr{
					Key:   "time",
					Value: slog.StringValue(a.Value.Time().Format("2006-01-02T15:04:05Z07:00")),
				}
			case slog.LevelKey:
				return slog.Attr{
					Key:   "level",
					Value: slog.StringValue(a.Value.String()),
				}
			case slog.MessageKey:
				return slog.Attr{
					Key:   "msg",
					Value: a.Value,
				}
			}
			return a
		},
	}
	var handler slog.Handler
	if cfg.Env == "local" {
		handler = slog.NewTextHandler(w, opts)
	} else {
		handler = slog.NewJSONHandler(w, opts)
	}

	base := slog.New(handler)

	return &SlogLogger{
		log: base.With(
			"app_name", cfg.AppName,
			"app_version", cfg.AppVersion,
		)}
}

func (l *SlogLogger) Info(msg string, fields ...any) {
	l.log.Info(msg, fields...)
}

func (l *SlogLogger) Debug(msg string, fields ...any) {
	l.log.Debug(msg, fields...)
}

func (l *SlogLogger) Error(msg string, fields ...any) {
	l.log.Error(msg, fields...)
}

func (l *SlogLogger) With(fields ...any) Logger {
	return &SlogLogger{
		log: l.log.With(fields...),
	}
}

func (c Config) logLevel() slog.Level {
	switch strings.ToLower(c.Level) {
	case "info":
		return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
