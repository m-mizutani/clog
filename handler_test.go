package clog_test

import (
	"bytes"
	"testing"

	"github.com/m-mizutani/clog"
	"github.com/m-mizutani/gt"
	"golang.org/x/exp/slog"
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
}
