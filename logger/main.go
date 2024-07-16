// Package logger provides a simple logging implementation to be used in conjunction with Puff.
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
	level slog.Level
}

func NewPuffSlogHandler(baseHandler slog.Handler, level slog.Level) *PuffSlogHandler {
	return &PuffSlogHandler{
		Handler: baseHandler,
		level:   level,
	}
}

func (h *PuffSlogHandler) Enabled(c context.Context, level slog.Level) bool {
	return level >= h.level
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

	return nil
}

func (h *PuffSlogHandler) SetLevel(level slog.Level) {
	h.level = level
}

func NewPuffLogger(c Config) *slog.Logger {
	var handler slog.Handler
	if c.UseJSON {
		handler =
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					AddSource: true,
					Level:     c.Level,
				},
			)
	} else {
		handler = slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: true,
				Level:     c.Level,
			},
		)
	}
	logger := slog.New(NewPuffSlogHandler(handler, c.Level))
	slog.SetDefault(logger)
	return logger
}

func DefaultPuffLogger() *slog.Logger {
	return NewPuffLogger(Config{
		UseJSON: true,
		Level:   slog.LevelInfo,
	})
}
