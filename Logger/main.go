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
	UseJSON       bool
	ColoredOutput bool
	*slog.HandlerOptions
}

type CustomSlogHandler struct {
	slog.Handler
}

func (h *CustomSlogHandler) Handle(c context.Context, r slog.Record) error {
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
	custom_options := slog.HandlerOptions{Level: c.Level, AddSource: c.AddSource}
	logger := slog.Default()
	if c.UseJSON {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &custom_options))
	}
	slog.SetDefault(logger)
	return logger
}

func DefaultLogger() *slog.Logger {
	custom_options := slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}

	h := CustomSlogHandler{Handler: slog.NewTextHandler(os.Stdout, &custom_options)}
	json_logger := slog.New(&h)
	slog.SetDefault(json_logger)
	return json_logger
}
