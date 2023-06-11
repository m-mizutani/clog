package clog

import (
	"fmt"
	"io"
	"strings"

	"github.com/k0kubun/pp/v3"
	"golang.org/x/exp/slog"
)

type AttrPrinter interface {
	Enter(group string)
	Exit(group string)
	Print(attr slog.Attr)
}

type basicPrinter struct {
	w      io.Writer
	cfg    *config
	groups []string
}

func newBasicPrinter(w io.Writer, cfg *config) basicPrinter {
	return basicPrinter{
		w:   w,
		cfg: cfg,
	}
}

func (x *basicPrinter) Enter(group string) {
	x.groups = append(x.groups, group)
}

func (x *basicPrinter) Exit(group string) {
	x.groups = x.groups[:len(x.groups)-1]
}

// LinearPrinter is a printer that prints attributes in a linear format.
func LinearPrinter(w io.Writer, cfg *config) AttrPrinter {
	return &linearPrinter{
		basicPrinter: newBasicPrinter(w, cfg),
	}
}

type linearPrinter struct {
	basicPrinter
}

func (x *linearPrinter) Print(attr slog.Attr) {
	var keyPrefix string
	if len(x.groups) > 0 {
		keyPrefix = strings.Join(x.groups, ".") + "."
	}
	key := keyPrefix + attr.Key

	if x.cfg.replaceAttr != nil {
		attr = x.cfg.replaceAttr(x.groups, attr)
	}

	p := fmt.Fprint
	if x.cfg.enableColor && x.cfg.colors.AttrKey != nil {
		p = x.cfg.colors.AttrKey.Fprint
	}
	p(x.w, key)

	p = fmt.Fprint
	p(x.w, "=")

	if x.cfg.enableColor && x.cfg.colors.AttrValue != nil {
		p = x.cfg.colors.AttrValue.Fprint
	}

	p(x.w, attr.Value.Resolve().String())
	p = fmt.Fprint
	p(x.w, " ")
}

// PrettyPrinter is a printer that prints attributes in a pretty format.
func PrettyPrinter(w io.Writer, cfg *config) AttrPrinter {
	p := &prettyPrinter{
		printer:      pp.New(),
		basicPrinter: newBasicPrinter(w, cfg),
	}
	p.printer.SetColoringEnabled(cfg.enableColor)
	return p
}

type prettyPrinter struct {
	printer *pp.PrettyPrinter
	basicPrinter
}

func (x *prettyPrinter) Print(attr slog.Attr) {
	var keyPrefix string
	if len(x.groups) > 0 {
		keyPrefix = strings.Join(x.groups, ".") + "."
	}
	key := keyPrefix + attr.Key

	if x.cfg.replaceAttr != nil {
		attr = x.cfg.replaceAttr(x.groups, attr)
	}

	p := fmt.Fprint
	p(x.w, "\n")

	if x.cfg.enableColor && x.cfg.colors.AttrKey != nil {
		p = x.cfg.colors.AttrKey.Fprint
	}
	p(x.w, key)

	p = fmt.Fprint
	p(x.w, " => ")

	if x.cfg.replaceAttr != nil {
		attr = x.cfg.replaceAttr(x.groups, attr)
	}
	x.printer.Fprint(x.w, attr.Value.Resolve().Any())
}

// IndentPrinter is a printer that prints attributes in a indented format.
func IndentPrinter(w io.Writer, cfg *config) AttrPrinter {
	return &indentPrinter{
		basicPrinter: newBasicPrinter(w, cfg),
	}
}

type indentPrinter struct {
	basicPrinter
}

func (x *indentPrinter) Enter(group string) {
	indent := "    " + strings.Repeat("  ", len(x.groups))
	fmt.Fprintf(x.w, "\n%s%s:", indent, group)
	x.basicPrinter.Enter(group)
}

func (x *indentPrinter) Print(attr slog.Attr) {
	indent := "    " + strings.Repeat("  ", len(x.groups))

	key := attr.Key
	if x.cfg.enableColor && x.cfg.colors.AttrKey != nil {
		key = x.cfg.colors.AttrKey.Sprint(key)
	}

	if x.cfg.replaceAttr != nil {
		attr = x.cfg.replaceAttr(x.groups, attr)
	}
	value := attr.Value.Resolve().String()
	if x.cfg.enableColor && x.cfg.colors.AttrValue != nil {
		value = x.cfg.colors.AttrValue.Sprint(value)
	}

	fmt.Fprintf(x.w, "\n%s%s: %s", indent, key, value)
}
