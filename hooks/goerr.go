package hooks

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/m-mizutani/clog"
	"github.com/m-mizutani/goerr/v2"
)

// goErrConfig holds configuration for GoErr hook.
type goErrConfig struct {
	withStackTrace bool
}

// GoErrOption is a functional option for GoErr hook.
type GoErrOption func(*goErrConfig)

// WithStackTrace enables or disables stack trace output.
// When enabled, the error message will include the full stack trace.
// Default is false (no stack trace).
func WithStackTrace(enable bool) GoErrOption {
	return func(cfg *goErrConfig) {
		cfg.withStackTrace = enable
	}
}

// GoErr creates an AttrHook for goerr.Error handling.
// It extracts the values from goerr.Error and prints the error message.
// Use WithStackTrace(true) to include stack trace in the output.
func GoErr(opts ...GoErrOption) clog.AttrHook {
	cfg := &goErrConfig{
		withStackTrace: false,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	return func(_ []string, attr slog.Attr) *clog.HandleAttr {
		goErr, ok := attr.Value.Any().(*goerr.Error)
		if !ok {
			return nil
		}

		var attrs []any
		for k, v := range goErr.Values() {
			attrs = append(attrs, slog.Any(k, v))
		}

		if cfg.withStackTrace {
			newAttr := slog.Group(attr.Key, attrs...)
			return &clog.HandleAttr{
				NewAttr: &newAttr,
				Defer: func(w io.Writer) {
					_, _ = fmt.Fprintf(w, "Error: %+v", goErr)
				},
			}
		}

		// Add error message to attributes only when stack trace is disabled
		attrs = append(attrs, slog.String("message", goErr.Error()))
		newAttr := slog.Group(attr.Key, attrs...)
		return &clog.HandleAttr{
			NewAttr: &newAttr,
		}
	}
}
