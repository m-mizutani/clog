package clog

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/exp/slog"
)

type AttrPrinter interface {
	Enter(group string)
	Exit(group string)
	Print(attr slog.Attr)
}

type linearPrinter struct {
	w      io.Writer
	groups []string
}

func LinearPrinter(w io.Writer) AttrPrinter {
	return &linearPrinter{
		w: w,
	}
}

func (x *linearPrinter) Enter(group string) {
	x.groups = append(x.groups, group)
}

func (x *linearPrinter) Exit(group string) {
	x.groups = x.groups[:len(x.groups)-1]
}

func (x *linearPrinter) Print(attr slog.Attr) {
	var keyPrefix string
	if len(x.groups) > 0 {
		keyPrefix = strings.Join(x.groups, ".") + "."
	}
	key := keyPrefix + attr.Key

	p := fmt.Fprint
	p(x.w, key)
	p(x.w, "=")

	p = color.New(color.FgHiWhite).Fprint
	p(x.w, attr.Value.String())
	p = fmt.Fprint
	p(x.w, " ")
}
