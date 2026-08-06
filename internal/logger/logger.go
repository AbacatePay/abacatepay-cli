package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync/atomic"

	"gopkg.in/natefinch/lumberjack.v2"
)

var consoleEnabled atomic.Bool

func init() {
	consoleEnabled.Store(true)
}

// ConsoleLevel derives the console handler's threshold from the configured
// file-log level. Every operational log call in this codebase (Info/Debug/
// Warn - there is no slog.Error call site; user-facing failures always go
// through output.Error/style.PrintError instead) is internal diagnostic
// detail, already reflected where it matters through the CLI's styled
// success/error/event output. Echoing it again as raw `level=INFO
// msg=...` text on stderr is redundant and visually inconsistent with that
// styling, so by default none of it reaches the console - only the file
// handler. Console only mirrors the file level once verbose mode (Debug) is
// on, when raw diagnostic detail is exactly what was asked for.
func ConsoleLevel(base slog.Level) slog.Level {
	if base <= slog.LevelDebug {
		return base
	}
	return slog.LevelError + 1
}

// SetConsoleEnabled toggles whether the console (stderr) log handler emits
// output. Callers that take over the terminal (e.g. a bubbletea alt-screen
// program) should disable it for the duration to avoid corrupting the
// display, then restore it afterwards. Log records still reach the file
// handler regardless of this setting.
func SetConsoleEnabled(enabled bool) {
	consoleEnabled.Store(enabled)
}

type toggleHandler struct {
	inner slog.Handler
}

func (h *toggleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return consoleEnabled.Load() && h.inner.Enabled(ctx, level)
}

func (h *toggleHandler) Handle(ctx context.Context, r slog.Record) error {
	if !consoleEnabled.Load() {
		return nil
	}
	return h.inner.Handle(ctx, r)
}

func (h *toggleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &toggleHandler{inner: h.inner.WithAttrs(attrs)}
}

func (h *toggleHandler) WithGroup(name string) slog.Handler {
	return &toggleHandler{inner: h.inner.WithGroup(name)}
}

type Config struct {
	LogDir     string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Level      slog.Level
}

func DefaultConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve home directory: %w", err)
	}

	return &Config{
		LogDir:     filepath.Join(homeDir, ".abacatepay", "logs"),
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
		Level:      slog.LevelInfo,
	}, nil
}

func Setup(cfg *Config) (*slog.Logger, error) {
	if err := os.MkdirAll(cfg.LogDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	logFile := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.LogDir, "abacatepay.log"),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	fileHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: cfg.Level})

	consoleHandler := &toggleHandler{inner: slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: ConsoleLevel(cfg.Level),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			return a
		},
	})}

	multiHandler := NewFanoutHandler(consoleHandler, fileHandler)

	logger := slog.New(multiHandler)
	slog.SetDefault(logger)

	return logger, nil
}

func NewTransactionLogger(cfg *Config) (*slog.Logger, error) {
	if err := os.MkdirAll(cfg.LogDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	logFile := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.LogDir, "transactions.log"),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	handler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(handler), nil
}

func NewConsoleLogger(level slog.Level) *slog.Logger {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler)
}

type FanoutHandler struct {
	handlers []slog.Handler
}

func NewFanoutHandler(handlers ...slog.Handler) *FanoutHandler {
	return &FanoutHandler{handlers: handlers}
}

func (h *FanoutHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *FanoutHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			_ = handler.Handle(ctx, r.Clone())
		}
	}
	return nil
}

func (h *FanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return NewFanoutHandler(handlers...)
}

func (h *FanoutHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return NewFanoutHandler(handlers...)
}
