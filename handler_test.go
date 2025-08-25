package clog_test

import (
	"bytes"
	"testing"

	"log/slog"

	"github.com/m-mizutani/clog"
	"github.com/m-mizutani/gt"
)

func TestWithGroup(t *testing.T) {
	w := &bytes.Buffer{}
	logger := slog.New(clog.New(
		clog.WithColor(false),
		clog.WithWriter(w),
	))
	logger.WithGroup("group1").Info("hello, world!", slog.String("foo", "bar"))

	gt.String(t, w.String()).
		Contains("INFO").
		Contains("hello, world!").
		Contains(`group1.foo="bar"`)
}

func TestGroup(t *testing.T) {
	w := &bytes.Buffer{}
	logger := slog.New(clog.New(
		clog.WithColor(false),
		clog.WithWriter(w),
	))
	logger.Info("hello, world!", slog.Group("group1", slog.String("foo", "bar")))

	gt.String(t, w.String()).
		Contains("INFO").
		Contains("hello, world!").
		Contains(`group1.foo="bar"`)
}

func TestGroupInGroup(t *testing.T) {
	testCases := map[string]struct {
		f func(l *slog.Logger)
	}{
		"record": {
			f: func(l *slog.Logger) {
				l.Info("hello, world!",
					slog.Group("group1",
						slog.Group("group2",
							slog.String("foo", "bar"),
						),
					),
				)
			},
		},
		"with": {
			f: func(l *slog.Logger) {
				l.WithGroup("group1").
					WithGroup("group2").
					Info("hello, world!", slog.String("foo", "bar"))
			},
		},
		"mix": {
			f: func(l *slog.Logger) {
				l.WithGroup("group1").Info("hello, world!",
					slog.Group("group2",
						slog.String("foo", "bar"),
					),
				)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			w := &bytes.Buffer{}
			logger := slog.New(clog.New(
				clog.WithColor(false),
				clog.WithWriter(w),
			))

			tc.f(logger)

			gt.String(t, w.String()).
				Contains("INFO").
				Contains("hello, world!").
				Contains(`group1.group2.foo="bar"`)
		})
	}
}

func TestWithAttrs(t *testing.T) {
	w := &bytes.Buffer{}
	logger := slog.New(clog.New(
		clog.WithColor(false),
		clog.WithWriter(w),
	))
	logger.
		With(slog.String("foo", "bar")).
		With(slog.String("hoge", "fuga")).
		Info("hello, world!")

	gt.String(t, w.String()).
		Contains("INFO").
		Contains("hello, world!").
		Contains(`foo="bar"`).
		Contains(`hoge="fuga"`)

	w.Reset()
	logger.
		With(slog.String("color", "red")).
		Info("good bye!")
	gt.String(t, w.String()).
		Contains("good bye!").
		Contains(`color="red"`).
		NotContains(`foo="bar"`)
}

func TestAttr(t *testing.T) {
	w := &bytes.Buffer{}
	logger := slog.New(clog.New(
		clog.WithColor(false),
		clog.WithWriter(w),
	))
	logger.Info("hello, world!", slog.String("foo", "bar"))
	gt.String(t, w.String()).Contains(`foo="bar"`)

	w.Reset()
	logger.Info("hello, again!", slog.String("hoge", "fuga"))
	gt.String(t, w.String()).
		Contains(`hoge="fuga"`).
		NotContains(`foo="bar"`)
}

// NOTE: This test is disabled for reducing unnecessary dependencies.
// If you need to test this feature, please get github.com/m-mizutani/masq and enable this test.
/*
type logV struct{}

func (v logV) LogValue() slog.Value {
	return slog.GroupValue(slog.String("v", "logV"))
}

func TestLogValuer(t *testing.T) {
	w := &bytes.Buffer{}
	logger := slog.New(clog.New(
		clog.WithColor(false),
		clog.WithWriter(w),
		clog.WithReplaceAttr(masq.New()),
	))
	logger.Info("hello, world!", slog.Any("g", logV{}))

	gt.String(t, w.String()).Contains(`v="logV"`)
}
*/

func TestHandlerWithLevelFormatter(t *testing.T) {
	// Test that custom level formatter is applied correctly
	customFormatter := func(level slog.Level) string {
		return "[" + level.String() + "]"
	}

	buf := &bytes.Buffer{}
	handler := clog.New(
		clog.WithWriter(buf),
		clog.WithLevelFormatter(customFormatter),
		clog.WithColor(false),
	)
	logger := slog.New(handler)

	logger.Info("test message", slog.String("key", "value"))
	output := buf.String()
	gt.S(t, output).Contains("[INFO]")
	gt.S(t, output).Contains("test message")
	gt.S(t, output).Contains(`key="value"`)
}

func TestHandlerWithLevelFormatterAndColor(t *testing.T) {
	// Test that level formatter works with color enabled
	customFormatter := func(level slog.Level) string {
		switch level {
		case slog.LevelInfo:
			return "INFO "
		case slog.LevelWarn:
			return "WARN "
		case slog.LevelError:
			return "ERROR"
		default:
			return level.String()
		}
	}

	buf := &bytes.Buffer{}
	handler := clog.New(
		clog.WithWriter(buf),
		clog.WithLevelFormatter(customFormatter),
		clog.WithColor(true), // Color enabled
	)
	logger := slog.New(handler)

	// Test different levels
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()
	// The formatted strings should be present (color codes will be added around them)
	gt.S(t, output).Contains("INFO ")
	gt.S(t, output).Contains("WARN ")
	gt.S(t, output).Contains("ERROR")
}

func TestHandlerWithLevelFormatterAndGroup(t *testing.T) {
	// Test that level formatter works with WithGroup
	customFormatter := func(level slog.Level) string {
		return "LEVEL:" + level.String()
	}

	buf := &bytes.Buffer{}
	handler := clog.New(
		clog.WithWriter(buf),
		clog.WithLevelFormatter(customFormatter),
		clog.WithColor(false),
	)
	logger := slog.New(handler)

	logger.WithGroup("mygroup").Info("grouped message", slog.String("key", "value"))

	output := buf.String()
	gt.S(t, output).Contains("LEVEL:INFO")
	gt.S(t, output).Contains("grouped message")
	gt.S(t, output).Contains(`mygroup.key="value"`)
}
