package clog

import (
	"fmt"
	"io"
	"strings"

	"log/slog"

	"github.com/k0kubun/pp/v3"
)

type AttrPrinter interface {
	Print(groups []string, attr slog.Attr)
}

type basicPrinter struct {
	w   io.Writer
	cfg *config
}

func newBasicPrinter(w io.Writer, cfg *config) basicPrinter {
	return basicPrinter{
		w:   w,
		cfg: cfg,
	}
}

func (x *basicPrinter) Enter(group string) {
}

func (x *basicPrinter) Exit(group string) {
}

func (x *basicPrinter) Defer() {
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

func (x *linearPrinter) Print(groups []string, attr slog.Attr) {
	if attr.Value.Kind() == slog.KindGroup {
		return
	}

	var keyPrefix string
	if len(groups) > 0 {
		keyPrefix = strings.Join(groups, ".") + "."
	}

	key := keyPrefix + attr.Key

	p := fmt.Fprint
	if x.cfg.enableColor && x.cfg.colors.AttrKey != nil {
		p = x.cfg.colors.AttrKey.Fprint
	}
	_, _ = p(x.w, key)

	p = fmt.Fprint
	_, _ = p(x.w, "=")

	if x.cfg.enableColor && x.cfg.colors.AttrValue != nil {
		p = x.cfg.colors.AttrValue.Fprint
	}

	value := valueToString(attr.Value)

	_, _ = p(x.w, value)
	p = fmt.Fprint
	_, _ = p(x.w, " ")
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

func (x *prettyPrinter) Print(groups []string, attr slog.Attr) {
	var keyPrefix string
	if len(groups) > 0 {
		keyPrefix = strings.Join(groups, ".") + "."
	}
	key := keyPrefix + attr.Key

	p := fmt.Fprint
	_, _ = p(x.w, "\n")

	if x.cfg.enableColor && x.cfg.colors.AttrKey != nil {
		p = x.cfg.colors.AttrKey.Fprint
	}
	_, _ = p(x.w, key)

	p = fmt.Fprint
	_, _ = p(x.w, " => ")
	_, _ = x.printer.Fprint(x.w, attr.Value.Any())
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

func (x *indentPrinter) Print(groups []string, attr slog.Attr) {
	if attr.Value.Kind() == slog.KindGroup {
		indent := strings.Repeat("  ", len(groups))
		_, _ = fmt.Fprintf(x.w, "\n%s%s:", indent, attr.Key)
		return
	}

	indent := strings.Repeat("  ", len(groups))

	key := attr.Key
	if x.cfg.enableColor && x.cfg.colors.AttrKey != nil {
		key = x.cfg.colors.AttrKey.Sprint(key)
	}

	if x.cfg.replaceAttr != nil {
		attr = x.cfg.replaceAttr(groups, attr)
	}

	value := valueToString(attr.Value.Resolve())
	if x.cfg.enableColor && x.cfg.colors.AttrValue != nil {
		value = x.cfg.colors.AttrValue.Sprint(value)
	}

	_, _ = fmt.Fprintf(x.w, "\n%s%s: %s", indent, key, value)
}

func valueToString(value slog.Value) string {
	switch value.Kind() {
	case slog.KindBool:
		return fmt.Sprintf("%v", value.Bool())
	case slog.KindString:
		return fmt.Sprintf("%q", value.String())
	case slog.KindTime:
		return fmt.Sprintf("%v", value.Time())
	case slog.KindDuration:
		return fmt.Sprintf("%v", value.Duration())
	case slog.KindAny:
		return fmt.Sprintf("%+v", value.Any())
	case slog.KindFloat64:
		return fmt.Sprintf("%v", value.Float64())
	case slog.KindInt64:
		return fmt.Sprintf("%v", value.Int64())
	case slog.KindUint64:
		return fmt.Sprintf("%v", value.Uint64())
	case slog.KindLogValuer:
		return value.LogValuer().LogValue().String()

	// Should not happen, but just in
	case slog.KindGroup:
		return fmt.Sprintf("%+v", value.Group())
	default:
		return fmt.Sprintf("%+v", value.Any())
	}
}
