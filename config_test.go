package clog_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"testing"

	"github.com/m-mizutani/clog"
	"github.com/m-mizutani/gt"
)

func TestWithLevelFormatter(t *testing.T) {
	// Test custom level formatter
	customFormatter := func(level slog.Level) string {
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

	buf := &bytes.Buffer{}
	handler := clog.New(
		clog.WithWriter(buf),
		clog.WithLevelFormatter(customFormatter),
		clog.WithColor(false),
	)
	logger := slog.New(handler)

	testCases := []struct {
		level    slog.Level
		logFunc  func(string, ...any)
		expected string
	}{
		{slog.LevelInfo, logger.Info, "INFO "},
		{slog.LevelWarn, logger.Warn, "WARN "},
		{slog.LevelError, logger.Error, "ERROR"},
	}

	for _, tc := range testCases {
		t.Run(tc.level.String(), func(t *testing.T) {
			buf.Reset()
			tc.logFunc("test message")
			output := buf.String()
			gt.S(t, output).Contains(tc.expected)
			gt.S(t, output).Contains("test message")
		})
	}
}

func TestWithLevelFormatterNil(t *testing.T) {
	// Test that nil formatter doesn't break the handler
	buf := &bytes.Buffer{}
	handler := clog.New(
		clog.WithWriter(buf),
		clog.WithLevelFormatter(nil), // Should use default
		clog.WithColor(false),
	)
	logger := slog.New(handler)

	logger.Info("test message")
	output := buf.String()
	gt.S(t, output).Contains("INFO")
	gt.S(t, output).Contains("test message")
}

func TestDefaultLevelFormatter(t *testing.T) {
	// Test that DefaultLevelFormatter works as expected
	testCases := []struct {
		level    slog.Level
		expected string
	}{
		{slog.LevelDebug, "DEBUG"},
		{slog.LevelInfo, "INFO"},
		{slog.LevelWarn, "WARN"},
		{slog.LevelError, "ERROR"},
	}

	for _, tc := range testCases {
		t.Run(tc.level.String(), func(t *testing.T) {
			result := clog.DefaultLevelFormatter(tc.level)
			gt.V(t, result).Equal(tc.expected)
		})
	}
}

func TestWithLevelFormatterDefault(t *testing.T) {
	// Test that without WithLevelFormatter, default behavior is preserved
	buf := &bytes.Buffer{}
	handler := clog.New(
		clog.WithWriter(buf),
		clog.WithColor(false),
	)
	logger := slog.New(handler)

	logger.Info("test message")
	output := buf.String()
	gt.S(t, output).Contains("INFO")
	gt.S(t, output).Contains("test message")
}
