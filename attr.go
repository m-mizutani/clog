package clog

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/m-mizutani/goerr/v2"
)

// HandleAttr is a struct that describes how to handle an attribute.
// NOTE: This feature is experimental and may be changed in the future.
// NOTE: This feature is available only in the LinearPrinter for now.
type HandleAttr struct {
	// NewAttr is a new attribute that replaces the original attribute. When this field is not nil, the original attribute is printed.
	NewAttr *slog.Attr

	// Defer is a function that is called after the all attributes are printed.
	Defer func(w io.Writer)
}

// AttrHook is a function that hooks attribute printing. When the function returns nil, the attribute is printed as usual. When the function returns a non-nil value, the attribute is handled according to the content of HandleAttr.
type AttrHook func(groups []string, attr slog.Attr) *HandleAttr

// GoerrHook is a hook function that hides the goerr.Error attribute and prints the error message.
//
// Deprecated: Use hooks.GoErr() instead. This function will be removed in a future version.
// The new hooks.GoErr() provides more options such as WithStackTrace to control stack trace output.
func GoerrHook(_ []string, attr slog.Attr) *HandleAttr {
	if goErr, ok := attr.Value.Any().(*goerr.Error); ok {
		var attrs []any
		for k, v := range goErr.Values() {
			attrs = append(attrs, slog.Any(k, v))
		}
		newAttr := slog.Group(attr.Key, attrs...)

		return &HandleAttr{
			NewAttr: &newAttr,
			Defer: func(w io.Writer) {
				_, _ = fmt.Fprintf(w, "Error: %+v", goErr)
			},
		}
	}

	return nil
}
