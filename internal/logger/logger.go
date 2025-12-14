package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

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
		return nil, fmt.Errorf("falha ao obter diretório home: %w", err)
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
	if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar diretório de logs: %w", err)
	}

	logFile := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.LogDir, "abacatepay.log"),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: cfg.Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger, nil
}

func NewTransactionLogger(cfg *Config) (*slog.Logger, error) {
	if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar diretório de logs: %w", err)
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
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler)
}
