package main

import (
	"fmt"
	"log/slog"

	"github.com/m-mizutani/clog"
)

func main() {
	fmt.Println("=== Default Level Formatter ===")
	demoDefault()

	fmt.Println("\n=== Custom Fixed Width Formatter (5 chars) ===")
	demoFixedWidth()

	fmt.Println("\n=== Custom Bracketed Formatter ===")
	demoBracketed()

	fmt.Println("\n=== Custom Uppercase Formatter ===")
	demoUppercase()
}

func demoDefault() {
	handler := clog.New(
		clog.WithColor(false),
		clog.WithSource(true),
	)
	logger := slog.New(handler)

	logger.Debug("Debug message", slog.String("key", "value"))
	logger.Info("Info message", slog.String("key", "value"))
	logger.Warn("Warning message", slog.String("key", "value"))
	logger.Error("Error message", slog.String("key", "value"))
}

func demoFixedWidth() {
	// Create a formatter that ensures all levels are exactly 5 characters
	fixedWidthFormatter := func(level slog.Level) string {
		switch level {
		case slog.LevelDebug:
			return "DEBUG"
		case slog.LevelInfo:
			return "INFO "
		case slog.LevelWarn:
			return "WARN "
		case slog.LevelError:
			return "ERROR"
		default:
			return fmt.Sprintf("%-5s", level.String())
		}
	}

	handler := clog.New(
		clog.WithLevelFormatter(fixedWidthFormatter),
		clog.WithColor(false),
		clog.WithSource(true),
	)
	logger := slog.New(handler)

	logger.Debug("Debug message", slog.String("key", "value"))
	logger.Info("Info message", slog.String("key", "value"))
	logger.Warn("Warning message", slog.String("key", "value"))
	logger.Error("Error message", slog.String("key", "value"))
}

func demoBracketed() {
	// Create a formatter that wraps levels in brackets
	bracketedFormatter := func(level slog.Level) string {
		return "[" + level.String() + "]"
	}

	handler := clog.New(
		clog.WithLevelFormatter(bracketedFormatter),
		clog.WithColor(false),
		clog.WithSource(true),
	)
	logger := slog.New(handler)

	logger.Info("Info message", slog.String("key", "value"))
	logger.Warn("Warning message", slog.String("key", "value"))
	logger.Error("Error message", slog.String("key", "value"))
}

func demoUppercase() {
	// Create a formatter that adds a prefix and uses uppercase
	uppercaseFormatter := func(level slog.Level) string {
		switch level {
		case slog.LevelDebug:
			return "DBG"
		case slog.LevelInfo:
			return "INF"
		case slog.LevelWarn:
			return "WRN"
		case slog.LevelError:
			return "ERR"
		default:
			return "UNK"
		}
	}

	handler := clog.New(
		clog.WithLevelFormatter(uppercaseFormatter),
		clog.WithColor(true), // Works with color too
		clog.WithSource(true),
	)
	logger := slog.New(handler)

	logger.Info("Info message", slog.String("key", "value"))
	logger.Warn("Warning message", slog.String("key", "value"))
	logger.Error("Error message", slog.String("key", "value"))

	// Show it works with groups too
	logger.WithGroup("app").Info("Grouped message", slog.Int("count", 42))
}
