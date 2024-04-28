package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

type Config struct {
	UseJSON bool
	Level   slog.Level
}

type PuffSlogHandler struct {
	slog.Handler
}

func (h *PuffSlogHandler) Handle(c context.Context, r slog.Record) error {
	level := fmt.Sprintf("%s:", r.Level.String())

	switch r.Level {
	case slog.LevelDebug:
		level = color.New(color.FgMagenta, color.Bold).Sprint(level)
	case slog.LevelInfo:
		level = color.New(color.FgBlue, color.Bold).Sprint(level)
	case slog.LevelWarn:
		level = color.New(color.FgYellow, color.Bold).Sprint(level)
	case slog.LevelError:
		level = color.New(color.FgRed, color.Bold).Sprint(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})
	attrs_formatted, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[2006-01-02 15:04:05]")
	msg := color.CyanString(r.Message)

	if len(fields) > 0 {
		fmt.Println(timeStr, level, msg, string(attrs_formatted))
	} else {
		fmt.Println(timeStr, level, msg)
	}

	return fmt.Errorf("Error in logging")
}

func InitLogger(c *Config) *slog.Logger {
	var logger *slog.Logger
	switch c.UseJSON {
	case true:
		logger = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					AddSource: true,
					Level:     c.Level,
				},
			),
		)
	case false:
		logger = DefaultLogger()
	}

	slog.SetDefault(logger)
	return logger
}

func DefaultLogger() *slog.Logger {
	opts := slog.HandlerOptions{
		AddSource: true,
	}
	h := PuffSlogHandler{Handler: slog.NewTextHandler(os.Stdout, &opts)}
	json_logger := slog.New(&h)
	slog.SetDefault(json_logger)
	return json_logger
}
