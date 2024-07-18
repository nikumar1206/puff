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
	config Config
}

func NewPuffSlogHandler(baseHandler slog.Handler, config Config) *PuffSlogHandler {
	return &PuffSlogHandler{
		Handler: baseHandler,
		config:  config,
	}
}

func (h *PuffSlogHandler) Enabled(c context.Context, level slog.Level) bool {
	return level >= h.config.Level
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
	// populate fields
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	timeStr := r.Time.Format("[2006-01-02 15:04:05]")
	if h.config.UseJSON {
		attrs_formatted, err := json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
		if len(fields) > 0 {
			fmt.Println(timeStr, level, r.Message, string(attrs_formatted))
		} else {
			fmt.Println(timeStr, level, r.Message)
		}
	} else {
		fmt.Println(timeStr, level, r.Message)
	}

	return nil
}

func (h *PuffSlogHandler) SetLevel(level slog.Level) {
	h.config.Level = level
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
	logger := slog.New(NewPuffSlogHandler(handler, c))
	slog.SetDefault(logger)
	return logger
}

func DefaultPuffLogger() *slog.Logger {
	return NewPuffLogger(Config{
		UseJSON: true,
		Level:   slog.LevelInfo,
	})
}
