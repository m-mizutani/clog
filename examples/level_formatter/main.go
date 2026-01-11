package main

import (
	"fmt"
	"log/slog"

	"github.com/m-mizutani/clog"
)

func main() {
	// Create a formatter that ensures all levels are exactly 5 characters
	// This works perfectly with color because the formatter is applied
	// before color codes are added
	fixedWidthFormatter := func(level slog.Level) string {
		switch level {
		case slog.LevelDebug:
			return "DEBUG" // 5 chars
		case slog.LevelInfo:
			return "INFO " // 5 chars with trailing space
		case slog.LevelWarn:
			return "WARN " // 5 chars with trailing space
		case slog.LevelError:
			return "ERROR" // 5 chars
		default:
			return fmt.Sprintf("%-5s", level.String())
		}
	}

	// Create handler with color enabled and fixed-width formatter
	handler := clog.New(
		clog.WithLevelFormatter(fixedWidthFormatter),
		clog.WithColor(true), // Color enabled
		clog.WithSource(true),
		clog.WithLevel(slog.LevelDebug), // Enable debug level
	)
	logger := slog.New(handler)

	fmt.Println("=== Fixed-width level formatting with color ===")
	fmt.Println("All levels are aligned at exactly 5 characters:")
	fmt.Println()

	// Demonstrate all log levels with fixed-width formatting
	logger.Debug("Debug message - lowest priority", slog.String("status", "verbose"))
	logger.Info("Info message - normal operations", slog.String("status", "ok"))
	logger.Warn("Warning message - be careful", slog.String("status", "warning"))
	logger.Error("Error message - something failed", slog.String("status", "error"))

	fmt.Println()
	fmt.Println("Works with grouped attributes too:")
	logger.WithGroup("app").Info("Application started", slog.Int("port", 8080))
	logger.WithGroup("db").Error("Database connection failed", slog.String("err", "timeout"))
	logger.WithGroup("api").Debug("API request received", slog.String("method", "GET"))
	logger.WithGroup("auth").Warn("Invalid token attempt", slog.String("ip", "192.168.1.1"))
}
