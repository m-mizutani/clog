package clog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	"log/slog"

	"github.com/m-mizutani/goerr/v2"
)

// Handler is a slog handler that writes logs to an io.Writer.
type Handler struct {
	cfg *config

	attrs []slog.Attr
	group string
	mutex *sync.Mutex

	parent *Handler
}

var _ slog.Handler = (*Handler)(nil)

// New creates a new handler.
func New(options ...Option) *Handler {
	h := &Handler{
		cfg:   newConfig(),
		mutex: &sync.Mutex{},
	}

	for _, option := range options {
		option(h.cfg)
	}

	return h
}

// clone returns a copy of the handler.
func (x *Handler) clone() *Handler {
	newHandler := &Handler{
		cfg:    x.cfg,
		parent: x,
		mutex:  x.mutex,
	}

	return newHandler
}

// Enabled implements slog.Handler.
func (x *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return x.cfg.level.Level() <= level
}

type stack struct {
	handlers []*Handler
}

func (x *stack) push(h *Handler) {
	x.handlers = append(x.handlers, h)
}

func (x *stack) pop() *Handler {
	if len(x.handlers) == 0 {
		return nil
	}

	h := x.handlers[len(x.handlers)-1]
	x.handlers = x.handlers[:len(x.handlers)-1]
	return h
}

// Handle implements slog.Handler.
func (x *Handler) Handle(ctx context.Context, record slog.Record) error {
	x = x.clone()
	buf := &bytes.Buffer{}

	log := &Log{
		logLevel:  record.Level,
		Timestamp: record.Time.Format(x.cfg.timeFmt),
		Elapsed:   elapsedDuration(),
		Level:     record.Level.String(),
		Message:   record.Message,
	}
	if record.Time.IsZero() {
		log.Timestamp = "(no time)"
	}

	if x.cfg.addSource && record.PC != 0 {
		src := getSource(record.PC)
		log.FileName = filepath.Base(src.FilePath)
		log.FilePath = src.FilePath
		log.FuncName = src.Func
		log.FileLine = src.Line
	}

	if x.cfg.enableColor {
		log = log.Coloring(x.cfg.colors)
	}

	if err := x.cfg.tmpl.Execute(buf, log); err != nil {
		return goerr.Wrap(err, "failed to execute template")
	}

	// print attrs
	record.Attrs(func(attr slog.Attr) bool {
		x.attrs = append(x.attrs, attr)
		return true
	})

	st := &stack{}
	for handler := x; handler != nil; handler = handler.parent {
		st.push(handler)
	}

	p := &printer{
		hooks: x.cfg.attrHooks,
		resolver: func(g []string, a slog.Attr) slog.Attr {
			newAttr := slog.Attr{
				Key:   a.Key,
				Value: a.Value.Resolve(),
			}
			if x.cfg.replaceAttr != nil && newAttr.Value.Kind() != slog.KindGroup {
				newAttr = x.cfg.replaceAttr(g, newAttr)
			}
			return newAttr
		},
		attrPrinter: x.cfg.newAttrPrinter(buf, x.cfg),
	}

	p.printStack(st)
	for i := len(p.defers) - 1; i >= 0; i-- {
		p.defers[i](buf)
	}

	fmt.Fprint(buf, "\n")

	x.mutex.Lock()
	defer x.mutex.Unlock()
	if _, err := x.cfg.w.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

type resolver func(groups []string, attr slog.Attr) slog.Attr

type printer struct {
	groups      []string
	hooks       []AttrHook
	defers      []func(w io.Writer)
	resolver    resolver
	attrPrinter AttrPrinter
}

func (x *printer) printStack(st *stack) {
	h := st.pop()
	if h == nil {
		return
	}

	if h.group != "" {
		x.groups = append(x.groups, h.group)
	}

	for _, attr := range h.attrs {
		x.printAttr(attr)
	}
	x.printStack(st)

	if h.group != "" {
		x.groups = x.groups[:len(x.groups)-1]
	}
}

func (x *printer) printAttr(attr slog.Attr) {
	if attr.Equal(slog.Attr{}) {
		return
	}

	if len(x.hooks) > 0 {
		for _, hook := range x.hooks {
			if handle := hook(x.groups, attr); handle != nil {
				if handle.Defer != nil {
					x.defers = append(x.defers, handle.Defer)
				}
				if handle.NewAttr != nil {
					attr = *handle.NewAttr
				}
			}
		}
	}

	attr = x.resolver(x.groups, attr)

	if slog.KindGroup == attr.Value.Kind() {
		x.groups = append(x.groups, attr.Key)
	}

	x.attrPrinter.Print(x.groups, attr)

	if slog.KindGroup == attr.Value.Kind() {
		for _, a := range attr.Value.Group() {
			x.printAttr(a)
		}
		x.groups = x.groups[:len(x.groups)-1]
	}
}

// WithAttrs implements slog.Handler.
func (x *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := x.clone()
	newHandler.attrs = attrs
	return newHandler
}

// WithGroup implements slog.Handler.
func (x *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return x
	}

	newHandler := x.clone()
	newHandler.group = name
	return newHandler
}
