package hooks_test

import (
	"bytes"
	"io"
	"log/slog"
	"testing"

	"github.com/m-mizutani/clog"
	"github.com/m-mizutani/clog/hooks"
	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/gt"
)

func toPtr[T any](v T) *T {
	return &v
}

func TestHandleAttr(t *testing.T) {
	type testCase struct {
		hook clog.AttrHook
		test func(t *testing.T, s string, attrs []slog.Attr)
	}

	runTest := func(tc testCase) func(t *testing.T) {
		return func(t *testing.T) {
			var hooked []slog.Attr
			var buf bytes.Buffer
			logger := slog.New(clog.New(
				clog.WithWriter(&buf),
				clog.WithAttrHook(func(groups []string, attr slog.Attr) *clog.HandleAttr {
					hooked = append(hooked, attr)
					return tc.hook(groups, attr)
				}),
			))
			logger.Info("hello, world!",
				slog.String("color", "blue"),
				slog.Any("number", 5),
				slog.Group("magic", slog.String("words", "timeless")),
			)

			tc.test(t, buf.String(), hooked)
		}
	}

	t.Run("no action", runTest(testCase{
		hook: func(_ []string, _ slog.Attr) *clog.HandleAttr {
			return nil
		},
		test: func(t *testing.T, s string, attrs []slog.Attr) {
			gt.S(t, s).
				Contains("hello, world!").
				Contains(`color="blue"`).
				Contains(`number=5`)
		},
	}))

	t.Run("replace attribute", runTest(testCase{
		hook: func(_ []string, attr slog.Attr) *clog.HandleAttr {
			if attr.Key == "color" {
				return &clog.HandleAttr{
					NewAttr: toPtr(slog.String("color", "red")),
				}
			}
			return nil
		},
		test: func(t *testing.T, s string, attrs []slog.Attr) {
			gt.S(t, s).
				Contains("hello, world!").
				Contains(`color="red"`).
				NotContains(`color="blue"`).
				Contains(`number=5`)
		},
	}))

	t.Run("defer action", runTest(testCase{
		hook: func(_ []string, attr slog.Attr) *clog.HandleAttr {
			if attr.Key == "color" {
				return &clog.HandleAttr{
					NewAttr: toPtr(slog.String("color", "red")),
					Defer: func(w io.Writer) {
						gt.R1(w.Write([]byte("deferred!"))).NoError(t)
					},
				}
			}
			return nil
		},
		test: func(t *testing.T, s string, attrs []slog.Attr) {
			gt.S(t, s).
				Contains("hello, world!").
				Contains(`color="red"`).
				Contains("deferred!").
				Contains(`number=5`)
		},
	}))
}

func TestGoErr(t *testing.T) {
	t.Run("without stack trace (default)", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(clog.New(
			clog.WithWriter(&buf),
			clog.WithAttrHook(hooks.GoErr()),
		))

		logger.Error("hello, world!", "err", goerr.New("something wrong", goerr.V("foo", "bar")))
		gt.S(t, buf.String()).
			NotContains("err.stacktrace=").
			NotContains("Error:").              // no deferred error output
			Contains(`message="something wrong"`). // message in attributes instead
			NotContains(".go:").                // stack trace should not appear
			Contains(`foo="bar"`)
	})

	t.Run("with stack trace enabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(clog.New(
			clog.WithWriter(&buf),
			clog.WithAttrHook(hooks.GoErr(hooks.WithStackTrace(true))),
		))

		logger.Error("hello, world!", "err", goerr.New("something wrong", goerr.V("foo", "bar")))
		gt.S(t, buf.String()).
			NotContains("err.stacktrace=").
			NotContains(`message=`).        // no message attribute when stack trace is enabled
			Contains("Error: something wrong").
			Contains(".go:"). // stack trace should appear
			Contains(`foo="bar"`)
	})

	t.Run("with stack trace explicitly disabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(clog.New(
			clog.WithWriter(&buf),
			clog.WithAttrHook(hooks.GoErr(hooks.WithStackTrace(false))),
		))

		logger.Error("hello, world!", "err", goerr.New("something wrong", goerr.V("foo", "bar")))
		gt.S(t, buf.String()).
			NotContains("err.stacktrace=").
			NotContains("Error:").              // no deferred error output
			Contains(`message="something wrong"`). // message in attributes instead
			NotContains(".go:").                // stack trace should not appear
			Contains(`foo="bar"`)
	})

	t.Run("non-goerr error returns nil", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(clog.New(
			clog.WithWriter(&buf),
			clog.WithAttrHook(hooks.GoErr()),
		))

		logger.Error("hello, world!", "err", "just a string")
		gt.S(t, buf.String()).
			Contains("hello, world!").
			Contains(`err="just a string"`)
	})
}
